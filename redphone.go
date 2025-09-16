package main

import (
	"crypto/rand"
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"

	"github.com/awnumar/memguard"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/chacha20poly1305"
)

// NATO encoding dictionary
var NATO = []string{
	"Zero-Zero......", "Zero-One.......", "Zero-Two.......", "Zero-Three.....", "Zero-Four......",
	"Zero-Five......", "Zero-Six.......", "Zero-Seven.....", "Zero-Eight.....", "Zero-Nine......",
	"Zero-Alfa......", "Zero-Bravo.....", "Zero-Charlie...", "Zero-Delta.....", "Zero-Echo......",
	"Zero-Foxtrot...", "One-Zero.......", "One-One........", "One-Two........", "One-Three......",
	"One-Four.......", "One-Five.......", "One-Six........", "One-Seven......", "One-Eight......",
	"One-Nine.......", "One-Alfa.......", "One-Bravo......", "One-Charlie....", "One-Delta......",
	"One-Echo.......", "One-Foxtrot....", "Two-Zero.......", "Two-One........", "Two-Two........",
	"Two-Three......", "Two-Four.......", "Two-Five.......", "Two-Six........", "Two-Seven......",
	"Two-Eight......", "Two-Nine.......", "Two-Alfa.......", "Two-Bravo......", "Two-Charlie....",
	"Two-Delta......", "Two-Echo.......", "Two-Foxtrot....", "Three-Zero.....", "Three-One......",
	"Three-Two......", "Three-Three....", "Three-Four.....", "Three-Five.....", "Three-Six......",
	"Three-Seven....", "Three-Eight....", "Three-Nine.....", "Three-Alfa.....", "Three-Bravo....",
	"Three-Charlie..", "Three-Delta....", "Three-Echo.....", "Three-Foxtrot..", "Four-Zero......",
	"Four-One.......", "Four-Two.......", "Four-Three.....", "Four-Four......", "Four-Five......",
	"Four-Six.......", "Four-Seven.....", "Four-Eight.....", "Four-Nine......", "Four-Alfa......",
	"Four-Bravo.....", "Four-Charlie...", "Four-Delta.....", "Four-Echo......", "Four-Foxtrot...",
	"Five-Zero......", "Five-One.......", "Five-Two.......", "Five-Three.....", "Five-Four......",
	"Five-Five......", "Five-Six.......", "Five-Seven.....", "Five-Eight.....", "Five-Nine......",
	"Five-Alfa......", "Five-Bravo.....", "Five-Charlie...", "Five-Delta.....", "Five-Echo......",
	"Five-Foxtrot...", "Six-Zero.......", "Six-One........", "Six-Two........", "Six-Three......",
	"Six-Four.......", "Six-Five.......", "Six-Six........", "Six-Seven......", "Six-Eight......",
	"Six-Nine.......", "Six-Alfa.......", "Six-Bravo......", "Six-Charlie....", "Six-Delta......",
	"Six-Echo.......", "Six-Foxtrot....", "Seven-Zero.....", "Seven-One......", "Seven-Two......",
	"Seven-Three....", "Seven-Four.....", "Seven-Five.....", "Seven-Six......", "Seven-Seven....",
	"Seven-Eight....", "Seven-Nine.....", "Seven-Alfa.....", "Seven-Bravo....", "Seven-Charlie..",
	"Seven-Delta....", "Seven-Echo.....", "Seven-Foxtrot..", "Eight-Zero.....", "Eight-One......",
	"Eight-Two......", "Eight-Three....", "Eight-Four.....", "Eight-Five.....", "Eight-Six......",
	"Eight-Seven....", "Eight-Eight....", "Eight-Nine.....", "Eight-Alfa.....", "Eight-Bravo....",
	"Eight-Charlie..", "Eight-Delta....", "Eight-Echo.....", "Eight-Foxtrot..", "Nine-Zero......",
	"Nine-One.......", "Nine-Two.......", "Nine-Three.....", "Nine-Four......", "Nine-Five......",
	"Nine-Six.......", "Nine-Seven.....", "Nine-Eight.....", "Nine-Nine......", "Nine-Alpha.....",
	"Nine-Bravo.....", "Nine-Charlie...", "Nine-Delta.....", "Nine-Echo......", "Nine-Foxtrot...",
	"Alfa-Zero......", "Alfa-One.......", "Alfa-Two.......", "Alfa-Three.....", "Alfa-Four......",
	"Alfa-Five......", "Alfa-Six.......", "Alfa-Seven.....", "Alfa-Eight.....", "Alfa-Nine......",
	"Alfa-Alfa......", "Alfa-Bravo.....", "Alfa-Charlie...", "Alfa-Delta.....", "Alfa-Echo......",
	"Alfa-Foxtrot...", "Bravo-Zero.....", "Bravo-One......", "Bravo-Two......", "Bravo-Three....",
	"Bravo-Four.....", "Bravo-Five.....", "Bravo-Six......", "Bravo-Seven....", "Bravo-Eight....",
	"Bravo-Nine.....", "Bravo-Alfa.....", "Bravo-Bravo....", "Bravo-Charlie..", "Bravo-Delta....",
	"Bravo-Echo.....", "Bravo-Foxtrot..", "Charlie-Zero...", "Charlie-One....", "Charlie-Two....",
	"Charlie-Three..", "Charlie-Four...", "Charlie-Five...", "Charlie-Six....", "Charlie-Seven..",
	"Charlie-Eight..", "Charlie-Nine...", "Charlie-Alfa...", "Charlie-Bravo..", "Charlie-Charlie",
	"Charlie-Delta..", "Charlie-Echo...", "Charlie-Foxtrot", "Delta-Zero.....", "Delta-One......",
	"Delta-Two......", "Delta-Three....", "Delta-Four.....", "Delta-Five.....", "Delta-Six......",
	"Delta-Seven....", "Delta-Eight....", "Delta-Nine.....", "Delta-Alfa.....", "Delta-Bravo....",
	"Delta-Charlie..", "Delta-Delta....", "Delta-Echo.....", "Delta-Foxtrot..", "Echo-Zero......",
	"Echo-One.......", "Echo-Two.......", "Echo-Three.....", "Echo-Four......", "Echo-Five......",
	"Echo-Six.......", "Echo-Seven.....", "Echo-Eight.....", "Echo-Nine......", "Echo-Alfa......",
	"Echo-Bravo.....", "Echo-Charlie...", "Echo-Delta.....", "Echo-Echo......", "Echo-Foxtrot...",
	"Foxtrot-Zero...", "Foxtrot-One....", "Foxtrot-Two....", "Foxtrot-Three..", "Foxtrot-Four...",
	"Foxtrot-Five...", "Foxtrot-Six....", "Foxtrot-Seven..", "Foxtrot-Eight..", "Foxtrot-Nine...",
	"Foxtrot-Alfa...", "Foxtrot-Bravo..", "Foxtrot-Charlie", "Foxtrot-Delta..", "Foxtrot-Echo...",
	"Foxtrot-Foxtrot",
}

// HEX encoding dictionary
var HEX = []string{
	"00", "01", "02", "03", "04", "05", "06", "07", "08", "09", "0A", "0B", "0C", "0D", "0E", "0F",
	"10", "11", "12", "13", "14", "15", "16", "17", "18", "19", "1A", "1B", "1C", "1D", "1E", "1F",
	"20", "21", "22", "23", "24", "25", "26", "27", "28", "29", "2A", "2B", "2C", "2D", "2E", "2F",
	"30", "31", "32", "33", "34", "35", "36", "37", "38", "39", "3A", "3B", "3C", "3D", "3E", "3F",
	"40", "41", "42", "43", "44", "45", "46", "47", "48", "49", "4A", "4B", "4C", "4D", "4E", "4F",
	"50", "51", "52", "53", "54", "55", "56", "57", "58", "59", "5A", "5B", "5C", "5D", "5E", "5F",
	"60", "61", "62", "63", "64", "65", "66", "67", "68", "69", "6A", "6B", "6C", "6D", "6E", "6F",
	"70", "71", "72", "73", "74", "75", "76", "77", "78", "79", "7A", "7B", "7C", "7D", "7E", "7F",
	"80", "81", "82", "83", "84", "85", "86", "87", "88", "89", "8A", "8B", "8C", "8D", "8E", "8F",
	"90", "91", "92", "93", "94", "95", "96", "97", "98", "99", "9A", "9B", "9C", "9D", "9E", "9F",
	"A0", "A1", "A2", "A3", "A4", "A5", "A6", "A7", "A8", "A9", "AA", "AB", "AC", "AD", "AE", "AF",
	"B0", "B1", "B2", "B3", "B4", "B5", "B6", "B7", "B8", "B9", "BA", "BB", "BC", "BD", "BE", "BF",
	"C0", "C1", "C2", "C3", "C4", "C5", "C6", "C7", "C8", "C9", "CA", "CB", "CC", "CD", "CE", "CF",
	"D0", "D1", "D2", "D3", "D4", "D5", "D6", "D7", "D8", "D9", "DA", "DB", "DC", "DD", "DE", "DF",
	"E0", "E1", "E2", "E3", "E4", "E5", "E6", "E7", "E8", "E9", "EA", "EB", "EC", "ED", "EE", "EF",
	"F0", "F1", "F2", "F3", "F4", "F5", "F6", "F7", "F8", "F9", "FA", "FB", "FC", "FD", "FE", "FF",
}

// encode: SMS encoding helper — converts 4-bit nibble to SMS character
func encode(input byte, secondChannel bool) rune {
	switch {
	case input <= 9 && secondChannel:
		return rune(input + 'A')
	case input <= 9:
		return rune(input + 'Q')
	case input >= 10 && input <= 15:
		return rune(input - 10 + 'K')
	default:
		panic("Invalid input")
	}
}

// decode: SMS decoding helper — converts SMS character back to 4-bit nibble
func decode(input rune, secondChannel bool) byte {
	switch {
	case input >= 'A' && input <= 'J' && secondChannel:
		return byte(input - 'A')
	case input >= 'Q' && input <= 'Z' && !secondChannel:
		return byte(input - 'Q')
	case input >= 'K' && input <= 'P':
		return byte(input-'K') + 10
	default:
		panic("Invalid input")
	}
}

// encodeBinary: Encodes binary data to SMS format (groups of 5 characters)
func encodeBinary(input []byte) string {
	var result strings.Builder
	secondChannel := false
	charactersInGroup := 0

	var allChars []rune
	for _, b := range input {
		high := b >> 4
		low := b & 0x0F
		encodedHigh := encode(high, secondChannel)
		encodedLow := encode(low, !secondChannel)
		allChars = append(allChars, encodedHigh, encodedLow)
		secondChannel = !secondChannel
	}

	for i, char := range allChars {
		result.WriteRune(char)
		charactersInGroup++
		if charactersInGroup == 5 {
			if i < len(allChars)-1 {
				result.WriteRune(' ')
			}
			charactersInGroup = 0
			if (i+1)%50 == 0 {
				result.WriteString("\n")
			}
		}
	}
	return result.String() + "\n"
}

// decodeBinary: Decodes SMS format back to binary data
func decodeBinary(input string) []byte {
	cleaned := strings.Map(func(r rune) rune {
		if r == ' ' || r == '\n' || r == '\r' {
			return -1
		}
		return r
	}, input)

	var result []byte
	secondChannel := false
	for i := 0; i < len(cleaned); i += 2 {
		if i+1 >= len(cleaned) {
			break
		}
		highChar, lowChar := rune(cleaned[i]), rune(cleaned[i+1])
		decodedHigh := decode(highChar, secondChannel)
		decodedLow := decode(lowChar, !secondChannel)
		result = append(result, decodedHigh<<4|decodedLow)
		secondChannel = !secondChannel
	}
	return result
}

// find: Helper to search for string in slice
func find(s string, dict []string) int {
	for i, v := range dict {
		if v == s {
			return i
		}
	}
	return -1
}

// Encrypt with random salt, nonce, and authentication (ChaCha20-Poly1305)
func encryptSecure(plaintext []byte, password string) ([]byte, error) {
	// 1. Generate random salt (16 bytes)
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}

	// 2. Derive key with Argon2id
	key := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, chacha20poly1305.KeySize)

	// 3. Create AEAD cipher
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	// 4. Generate random nonce (24 bytes for X variant)
	nonce := make([]byte, aead.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// 5. Encrypt + authenticate
	ciphertext := aead.Seal(nil, nonce, plaintext, nil)

	// 6. Prepend salt + nonce to ciphertext
	result := make([]byte, 0, 16+24+len(ciphertext))
	result = append(result, salt...)
	result = append(result, nonce...)
	result = append(result, ciphertext...)

	return result, nil
}

// Decrypt with authentication — detects tampering
func decryptSecure(encrypted []byte, password string) ([]byte, error) {
	if len(encrypted) < 16+24+16 { // salt + nonce + min tag
		return nil, fmt.Errorf("ciphertext too short — possible manipulation or corruption")
	}

	salt := encrypted[:16]
	nonce := encrypted[16:40]
	ciphertext := encrypted[40:]

	key := argon2.IDKey([]byte(password), salt, 3, 64*1024, 4, chacha20poly1305.KeySize)

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryption failed — message may be corrupted, manipulated, or password is wrong")
	}

	return plaintext, nil
}

// encodeNATO: Encodes binary data to NATO phonetic alphabet with line numbers
func encodeNATO(data []byte) string {
	var words []string
	for _, b := range data {
		if int(b) < len(NATO) {
			words = append(words, NATO[b])
		} else {
			words = append(words, "INVALID")
		}
	}

	var result strings.Builder
	lineNumber := 1
	for i := 0; i < len(words); i += 5 {
		end := i + 5
		if end > len(words) {
			end = len(words)
		}
		if i > 0 {
			result.WriteString("\r\n")
		}
		result.WriteString(fmt.Sprintf("%d\t%s", lineNumber, strings.Join(words[i:end], " ")))
		lineNumber++
	}
	return result.String()
}

// decodeHEX: Decodes HEX input (with line numbers) back to binary
func decodeHEX(text string) ([]byte, error) {
	if strings.TrimSpace(text) == "" {
		return nil, fmt.Errorf("input is empty")
	}

	lines := strings.Split(text, "\r\n")
	var hexStrings []string
	for _, line := range lines {
		parts := strings.SplitN(line, "\t", 2)
		if len(parts) > 1 {
			line = parts[1] // Remove line number
		}
		hexStrings = append(hexStrings, strings.Fields(line)...)
	}

	var result []byte
	for _, hexStr := range hexStrings {
		index := find(hexStr, HEX)
		if index == -1 {
			return nil, fmt.Errorf("hex string not found in dictionary: %s", hexStr)
		}
		result = append(result, byte(index))
	}

	return result, nil
}

// secureBytes: Wraps byte slice in secure memory using memguard
func secureBytes(data []byte) *memguard.LockedBuffer {
	return memguard.NewBufferFromBytes(data)
}

// showPasswordDialog: Clean dialog — no hint text, enforce 12+ chars
func showPasswordDialog(title string, callback func(string), parent fyne.Window) {
	passwordEntry := widget.NewPasswordEntry()
	dialogContent := container.NewVBox(
		passwordEntry,
	)

	d := dialog.NewCustomConfirm(title, "OK", "Cancel", dialogContent, func(confirmed bool) {
		if confirmed {
			if passwordEntry.Text == "" {
				dialog.ShowError(fmt.Errorf("password cannot be empty"), parent)
				return
			}
			if len(passwordEntry.Text) < 12 {
				dialog.ShowError(fmt.Errorf("password must be at least 12 characters long"), parent)
				return
			}
			callback(passwordEntry.Text)
		}
	}, parent)

	d.Resize(fyne.NewSize(400, 150))
	d.Show()
}

func main() {
	a := app.New()
	w := a.NewWindow("Red Phone")
	isDarkTheme := true
	a.Settings().SetTheme(theme.DarkTheme())
	isVoiceMode := false

	// Theme toggle
	var themeToggle *widget.Button
	themeToggle = widget.NewButtonWithIcon("", theme.ViewRefreshIcon(), func() {
		if isDarkTheme {
			a.Settings().SetTheme(theme.LightTheme())
			isDarkTheme = false
		} else {
			a.Settings().SetTheme(theme.DarkTheme())
			isDarkTheme = true
		}
	})
	themeToggle.Importance = widget.LowImportance

	// Voice/SMS toggle
	var voiceToggle *widget.Button
	voiceToggle = widget.NewButton("SMS", func() {
		isVoiceMode = !isVoiceMode
		if isVoiceMode {
			voiceToggle.SetText("Voice")
		} else {
			voiceToggle.SetText("SMS")
		}
	})

	// Text area
	textArea := widget.NewMultiLineEntry()
	textArea.Wrapping = fyne.TextWrapBreak
	textArea.TextStyle.Monospace = true
	textArea.PlaceHolder = "Enter text..."
	textAreaScroll := container.NewScroll(textArea)
	textAreaScroll.SetMinSize(fyne.NewSize(780, 400))

	// Status bar
	statusBar := widget.NewLabel("Ready")
	statusBar.TextStyle.Monospace = true
	statusBar.TextStyle.Italic = true

	// Encrypt button
	encryptBtn := widget.NewButtonWithIcon("Encrypt", theme.MailComposeIcon(), func() {
		if strings.TrimSpace(textArea.Text) == "" {
			statusBar.SetText("Error: Input field is empty")
			dialog.ShowError(fmt.Errorf("input field cannot be empty"), w)
			return
		}

		showPasswordDialog("Enter Password", func(password string) {
			plaintext := []byte(textArea.Text)

			encrypted, err := encryptSecure(plaintext, password)
			if err != nil {
				statusBar.SetText("Encryption failed: " + err.Error())
				return
			}

			var result string
			if isVoiceMode {
				result = encodeNATO(encrypted)
				statusBar.SetText("Encrypted")
			} else {
				result = encodeBinary(encrypted)
				statusBar.SetText("Encrypted")
			}

			textArea.SetText(result)
		}, w)
	})

	// Decrypt button
	decryptBtn := widget.NewButtonWithIcon("Decrypt", theme.MailForwardIcon(), func() {
		if strings.TrimSpace(textArea.Text) == "" {
			statusBar.SetText("Error: Input field is empty")
			dialog.ShowError(fmt.Errorf("input field cannot be empty"), w)
			return
		}

		showPasswordDialog("Enter Password", func(password string) {
			var decoded []byte
			var decodeErr error

			if isVoiceMode {
				decoded, decodeErr = decodeHEX(textArea.Text)
				if decodeErr != nil {
					statusBar.SetText("HEX decode failed: " + decodeErr.Error())
					return
				}
			} else {
				decoded = decodeBinary(textArea.Text)
			}

			plaintext, err := decryptSecure(decoded, password)
			if err != nil {
				statusBar.SetText("Decryption failed: " + err.Error())
				return
			}

			textArea.SetText(string(plaintext))
			statusBar.SetText("Decrypted — Authentication verified, message integrity confirmed")
		}, w)
	})

	// Clear button
	clearBtn := widget.NewButtonWithIcon("Clear", theme.DeleteIcon(), func() {
		textArea.SetText("")
		app := fyne.CurrentApp()
		windows := app.Driver().AllWindows()
		if len(windows) > 0 {
			windows[0].Clipboard().SetContent("")
			statusBar.SetText("All data cleared from display and clipboard")
		} else {
			statusBar.SetText("Display cleared (clipboard access failed)")
		}
	})

	// UI Layout
	header := container.NewHBox(
		layout.NewSpacer(),
		voiceToggle,
		themeToggle,
	)

	buttonContainer := container.NewHBox(
		encryptBtn,
		decryptBtn,
		clearBtn,
	)

	bottomPanel := container.NewVBox(
		container.NewCenter(buttonContainer),
		widget.NewSeparator(),
		container.NewHBox(
			statusBar,
			layout.NewSpacer(),
		),
	)

	mainContent := container.NewBorder(
		header,
		bottomPanel,
		nil,
		nil,
		textAreaScroll,
	)

	w.SetContent(mainContent)
	w.Resize(fyne.NewSize(800, 600))
	w.CenterOnScreen()
	w.SetTitle("Red Phone")

	w.ShowAndRun()
}
