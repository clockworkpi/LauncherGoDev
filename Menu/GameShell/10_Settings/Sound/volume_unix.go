// +build !windows
/*
 * Copied from https://github.com/itchyny/volume-go,  MIT License
 */
package Sound

import (
	"errors"
	//"fmt"
	//"os"
	//"os/exec"
	//"strings"
	"github.com/clockworkpi/LauncherGoDev/sysgo/UI"
)

// GetVolume returns the current volume (0 to 100).
func GetVolume() (int, error) {
	out, err := UI.ExecCmd(getVolumeCmd())
	if err != nil {
		return 0, err
	}
	return parseVolume(string(out))
}

// SetVolume sets the sound volume to the specified value.
func SetVolume(volume int) error {
	if volume < 0 || 100 < volume {
		return errors.New("out of valid volume range")
	}
	_, err := UI.ExecCmd(setVolumeCmd(volume))
	return err
}

// IncreaseVolume increases (or decreases) the audio volume by the specified value.
func IncreaseVolume(diff int) error {
	_, err := UI.ExecCmd(increaseVolumeCmd(diff))
	return err
}

// GetMuted returns the current muted status.
func GetMuted() (bool, error) {
	out, err := UI.ExecCmd(getMutedCmd())
	if err != nil {
		return false, err
	}
	return parseMuted(string(out))
}

// Mute mutes the audio.
func Mute() error {
	_, err := UI.ExecCmd(muteCmd())
	return err
}

// Unmute unmutes the audio.
func Unmute() error {
	_, err := UI.ExecCmd(unmuteCmd())
	return err
}
