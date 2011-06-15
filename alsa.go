// alsa package is the simple wrapper for C alsa binding library.
package alsa

// #include <alsa/asoundlib.h>
import "C"

import (
	"os"
	"unsafe"
	"fmt"
)

// Alsa stream type. Playback or capture.
type StreamType C.snd_pcm_stream_t

// Stream type constants.
const (
	// Playback stream
	StreamTypePlayback = C.SND_PCM_STREAM_PLAYBACK
	// Capture stream
	StreamTypeCapture = C.SND_PCM_STREAM_CAPTURE
)

// Sample type.
type SampleFormat C.snd_pcm_format_t

const (
	// Unknown
	SampleFormatUnknown = C.SND_PCM_FORMAT_UNKNOWN
	// Signed 8 bit
	SampleFormatS8 = C.SND_PCM_FORMAT_S8
	// Unsigned 8 bit
	SampleFormatU8 = C.SND_PCM_FORMAT_U8
	// Signed 16 bit Little Endian
	SampleFormatS16LE = C.SND_PCM_FORMAT_S16_LE
	// Signed 16 bit Big Endian
	SampleFromatS16BE = C.SND_PCM_FORMAT_S16_BE
	// Unsigned 16 bit Little Endian
	SampleFormatU16LE = C.SND_PCM_FORMAT_U16_LE
	// Unsigned 16 bit Big Endian
	SampleFormatU16BE = C.SND_PCM_FORMAT_U16_BE
	/*
	 SND_PCM_FORMAT_S24_LE 	Signed 24 bit Little Endian using low three bytes in 32-bit word
	 SND_PCM_FORMAT_S24_BE 	Signed 24 bit Big Endian using low three bytes in 32-bit word
	 SND_PCM_FORMAT_U24_LE 	Unsigned 24 bit Little Endian using low three bytes in 32-bit word
	 SND_PCM_FORMAT_U24_BE 	Unsigned 24 bit Big Endian using low three bytes in 32-bit word
	 SND_PCM_FORMAT_S32_LE 	Signed 32 bit Little Endian
	 SND_PCM_FORMAT_S32_BE 	Signed 32 bit Big Endian
	 SND_PCM_FORMAT_U32_LE 	Unsigned 32 bit Little Endian
	 SND_PCM_FORMAT_U32_BE 	Unsigned 32 bit Big Endian
	 SND_PCM_FORMAT_FLOAT_LE 	Float 32 bit Little Endian, Range -1.0 to 1.0
	 SND_PCM_FORMAT_FLOAT_BE 	Float 32 bit Big Endian, Range -1.0 to 1.0
	 SND_PCM_FORMAT_FLOAT64_LE 	Float 64 bit Little Endian, Range -1.0 to 1.0
	 SND_PCM_FORMAT_FLOAT64_BE 	Float 64 bit Big Endian, Range -1.0 to 1.0
	 SND_PCM_FORMAT_IEC958_SUBFRAME_LE 	IEC-958 Little Endian
	 SND_PCM_FORMAT_IEC958_SUBFRAME_BE 	IEC-958 Big Endian
	 SND_PCM_FORMAT_MU_LAW 	Mu-Law
	 SND_PCM_FORMAT_A_LAW 	A-Law
	 SND_PCM_FORMAT_IMA_ADPCM 	Ima-ADPCM
	 SND_PCM_FORMAT_MPEG 	MPEG
	 SND_PCM_FORMAT_GSM 	GSM
	 SND_PCM_FORMAT_SPECIAL 	Special
	 SND_PCM_FORMAT_S24_3LE 	Signed 24bit Little Endian in 3bytes format
	 SND_PCM_FORMAT_S24_3BE 	Signed 24bit Big Endian in 3bytes format
	 SND_PCM_FORMAT_U24_3LE 	Unsigned 24bit Little Endian in 3bytes format
	 SND_PCM_FORMAT_U24_3BE 	Unsigned 24bit Big Endian in 3bytes format
	 SND_PCM_FORMAT_S20_3LE 	Signed 20bit Little Endian in 3bytes format
	 SND_PCM_FORMAT_S20_3BE 	Signed 20bit Big Endian in 3bytes format
	 SND_PCM_FORMAT_U20_3LE 	Unsigned 20bit Little Endian in 3bytes format
	 SND_PCM_FORMAT_U20_3BE 	Unsigned 20bit Big Endian in 3bytes format
	 SND_PCM_FORMAT_S18_3LE 	Signed 18bit Little Endian in 3bytes format
	 SND_PCM_FORMAT_S18_3BE 	Signed 18bit Big Endian in 3bytes format
	 SND_PCM_FORMAT_U18_3LE 	Unsigned 18bit Little Endian in 3bytes format
	 SND_PCM_FORMAT_U18_3BE 	Unsigned 18bit Big Endian in 3bytes format
	 SND_PCM_FORMAT_S16 	Signed 16 bit CPU endian
	 SND_PCM_FORMAT_U16 	Unsigned 16 bit CPU endian
	 SND_PCM_FORMAT_S24 	Signed 24 bit CPU endian
	 SND_PCM_FORMAT_U24 	Unsigned 24 bit CPU endian
	 SND_PCM_FORMAT_S32 	Signed 32 bit CPU endian
	 SND_PCM_FORMAT_U32 	Unsigned 32 bit CPU endian
	 SND_PCM_FORMAT_FLOAT 	Float 32 bit CPU endian
	 SND_PCM_FORMAT_FLOAT64 	Float 64 bit CPU endian
	 SND_PCM_FORMAT_IEC958_SUBFRAME 	IEC-958 CPU Endian 
	*/
)

// Open mode constants.
const (
	ModeBlock    = 0
	ModeNonblock = C.SND_PCM_NONBLOCK
	ModeAsync    = C.SND_PCM_ASYNC
)

// Handle represents ALSA stream handler.
type Handle struct {
	cHandle *C.snd_pcm_t
	// Used samples format (size, endianness, signed).
	SampleFormat SampleFormat
	// Sample rate in Hz. Usual 44100.
	SampleRate int
	// Channels in the stream. 2 for stereo.
	Channels int
}

// New returns newly initialized ALSA handler.
func New() *Handle {
	handler := new(Handle)

	return handler
}

// Open opens a stream.
func (handle *Handle) Open(device string, streamType StreamType, mode int) os.Error {
	cDevice := C.CString(device)
	defer C.free(unsafe.Pointer(cDevice))

	err := C.snd_pcm_open(&(handle.cHandle), cDevice,
		_Ctypedef_snd_pcm_stream_t(streamType),
		_Ctype_int(mode))

	if err < 0 {
		return os.NewError(fmt.Sprintf("Cannot open audio device '%s'. %s",
			device, strError(err)))
	}

	return nil
}

// ApplyHwParams changes ALSA hardware parameters for the current stream.
func (handle *Handle) ApplyHwParams() os.Error {
	var cHwParams *C.snd_pcm_hw_params_t

	err := C.snd_pcm_hw_params_malloc(&cHwParams)
	if err < 0 {
		return os.NewError(fmt.Sprintf("Cannot allocate hardware parameter structure. %s",
			strError(err)))
	}

	err = C.snd_pcm_hw_params_any(handle.cHandle, cHwParams)
	if err < 0 {
		return os.NewError(fmt.Sprintf("Cannot initialize hardware parameter structure. %s",
			strError(err)))
	}

	err = C.snd_pcm_hw_params_set_access(handle.cHandle, cHwParams, C.SND_PCM_ACCESS_RW_INTERLEAVED)
	if err < 0 {
		return os.NewError(fmt.Sprintf("Cannot set access type. %s",
			strError(err)))
	}

	err = C.snd_pcm_hw_params_set_format(handle.cHandle, cHwParams, _Ctypedef_snd_pcm_format_t(handle.SampleFormat))
	if err < 0 {
		return os.NewError(fmt.Sprintf("Cannot set sample format. %s",
			strError(err)))
	}

	var cSampleRate _Ctype_uint = _Ctype_uint(handle.SampleRate)
	err = C.snd_pcm_hw_params_set_rate_near(handle.cHandle, cHwParams, &cSampleRate, nil)
	if err < 0 {
		return os.NewError(fmt.Sprintf("Cannot set sample rate. %s",
			strError(err)))
	}

	err = C.snd_pcm_hw_params_set_channels(handle.cHandle, cHwParams, _Ctype_uint(handle.Channels))
	if err < 0 {
		return os.NewError(fmt.Sprintf("Cannot set number of channels. %s",
			strError(err)))
	}

	// Drain current data and make sure we aren't underrun.
	C.snd_pcm_drain(handle.cHandle)

	err = C.snd_pcm_hw_params(handle.cHandle, cHwParams)
	if err < 0 {
		return os.NewError(fmt.Sprintf("Cannot set hardware parameters. %s",
			strError(err)))
	}

	C.snd_pcm_hw_params_free(cHwParams)

	return nil
}

// Wait waits till buffer will be free for some new portion of data or
// delay time is runs out.
func (handle *Handle) Wait(maxDelay int) os.Error {
	err := C.snd_pcm_wait(handle.cHandle, _Ctype_int(maxDelay))
	if err < 0 {
		return os.NewError(fmt.Sprintf("Pool failed. %s", strError(err)))
	}

	return nil
}

// AvailUpdate returns number of bytes ready to be read/written.
func (handle *Handle) AvailUpdate() (freeBytes int, err os.Error) {
	frames := C.snd_pcm_avail_update(handle.cHandle)
	if frames < 0 {
		return 0, os.NewError(fmt.Sprintf("Retriving free buffer size failed. %s", strError(_Ctype_int(frames))))
	}

	return int(frames) * handle.FrameSize(), nil
}

// Write writes given PCM data.
// Returns wrote value is total bytes was written.
func (handle *Handle) Write(buf []byte) (wrote int, err os.Error) {
	frames := len(buf) / handle.SampleSize() / handle.Channels
	wrote = int(C.snd_pcm_writei(handle.cHandle, unsafe.Pointer(&buf[0]), _Ctypedef_snd_pcm_uframes_t(frames)))
	if wrote < 0 {
		return 0, os.NewError(fmt.Sprintf("Write failed. %s", strError(_Ctype_int(wrote))))
	}

	wrote *= handle.FrameSize()

	return wrote, nil
}

// Close closes stream and release the handler.
func (handle *Handle) Close() {
	C.snd_pcm_close(handle.cHandle)
}

// SampleSize returns one sample size in bytes.
func (handle *Handle) SampleSize() int {
	switch handle.SampleFormat {
	case SampleFormatS8, SampleFormatU8:
		return 1
	case SampleFormatS16LE, SampleFromatS16BE,
		SampleFormatU16LE, SampleFormatU16BE:
		return 2
	}

	return 1
}

// FrameSize returns size of one frame in bytes.
func (handle *Handle) FrameSize() int {
	return handle.SampleSize() * handle.Channels
}

// strError retruns string description of ALSA error by its code.
func strError(err _Ctype_int) string {
	cErrMsg := C.snd_strerror(err)

	return C.GoString(cErrMsg)
}
