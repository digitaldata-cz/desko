package desko

import (
	"testing"

	"github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/assert"
	// "github.com/stretchr/testify/mock"
)

/*
// MockHIDDevice is a mock implementation of the hid.Device interface
type MockHIDDevice struct {
	mock.Mock
}

func (m *MockHIDDevice) Write(b []byte) (int, error) {
	args := m.Called(b)
	return args.Int(0), args.Error(1)
}

func (m *MockHIDDevice) Read(b []byte) (int, error) {
	args := m.Called(b)
	copy(b, args.Get(0).([]byte))
	return args.Int(1), args.Error(2)
}

func (m *MockHIDDevice) Close() error {
	args := m.Called()
	return args.Error(0)
}

*/

/*
func TestStartReading(t *testing.T) {
	mockDevice := new(MockHIDDevice)
	mockDevice.On("Write", []byte{0x20, 0x00}).Return(2, nil)
	mockDevice.On("Write", []byte{0x30, 0x00}).Return(2, nil)
	//	mockDevice.On("Read", mock.Anything).Return([]byte{0x1c, 0x02, 'P', '<', 'G', 'B', 'R', 'J', 'O', 'H', 'N', 'S', 'O', 'N', '<', '<', 'J', 'O', 'H', 'N', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<')
	return nil
}
*/

func TestParseICAO(t *testing.T) {
	tests := []struct {
		name     string
		input    ReaderData
		expected IcaoDocument
	}{
		{
			name: "Valid ICAO 9303 - Passport",
			input: ReaderData{
				Type: "MRZ",
				RawData: [][]byte{
					[]byte("P<GBRJOHNSON<<JOHN<<<<<<<<<<<<<<<<<<<<<<<<<<"),
					[]byte("1234567897GBR8001014M2501017<<<<<<<<<<<<<<02"),
				},
			},
			expected: IcaoDocument{
				IcaoType:         "P",
				IcaoSubtype:      "",
				Country:          "GBR",
				Number:           "123456789",
				NumberChecksum:   "7",
				NumberChecksumOk: true,
				Name:             "JOHN",
				Surname:          "JOHNSON",
				Sex:              "M",
				Nationality:      "GBR",
				Birth: struct {
					Year       string `json:"year" xml:"year"`
					Month      string `json:"month" xml:"month"`
					Day        string `json:"day" xml:"day"`
					Checksum   string `json:"checksum" xml:"checksum"`
					ChecksumOk bool   `json:"checksumok" xml:"checksumok"`
				}{
					Year:       "80",
					Month:      "01",
					Day:        "01",
					Checksum:   "4",
					ChecksumOk: true,
				},
				Expire: struct {
					Year       string `json:"year" xml:"year"`
					Month      string `json:"month" xml:"month"`
					Day        string `json:"day" xml:"day"`
					Checksum   string `json:"checksum" xml:"checksum"`
					ChecksumOk bool   `json:"checksumok" xml:"checksumok"`
				}{
					Year:       "25",
					Month:      "01",
					Day:        "01",
					Checksum:   "7",
					ChecksumOk: true,
				},
			},
		},
		{
			name: "Valid ICAO 9303 - ID Card",
			input: ReaderData{
				Type: "MRZ",
				RawData: [][]byte{
					[]byte("IDCZE9980013435<<<<<<<<<<<<<<<"),
					[]byte("8110088F2201100CZE<<<<<<<<<<<6"),
					[]byte("SPECIMEN<<VZOR<<<<<<<<<<<<<<<<"),
				},
			},
			expected: IcaoDocument{
				IcaoType:         "I",
				IcaoSubtype:      "D",
				Country:          "CZE",
				Number:           "998001343",
				NumberChecksum:   "5",
				NumberChecksumOk: true,
				Name:             "VZOR",
				Surname:          "SPECIMEN",
				Sex:              "F",
				Nationality:      "CZE",
				Birth: struct {
					Year       string `json:"year" xml:"year"`
					Month      string `json:"month" xml:"month"`
					Day        string `json:"day" xml:"day"`
					Checksum   string `json:"checksum" xml:"checksum"`
					ChecksumOk bool   `json:"checksumok" xml:"checksumok"`
				}{
					Year:       "81",
					Month:      "10",
					Day:        "08",
					Checksum:   "8",
					ChecksumOk: true,
				},
				Expire: struct {
					Year       string `json:"year" xml:"year"`
					Month      string `json:"month" xml:"month"`
					Day        string `json:"day" xml:"day"`
					Checksum   string `json:"checksum" xml:"checksum"`
					ChecksumOk bool   `json:"checksumok" xml:"checksumok"`
				}{
					Year:       "22",
					Month:      "01",
					Day:        "10",
					Checksum:   "0",
					ChecksumOk: true,
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ParseICAO(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

/*
func TestStartReading(t *testing.T) {
	mockDevice := new(MockHIDDevice)
	mockDevice.On("Write", []byte{0x20, 0x00}).Return(2, nil)
	mockDevice.On("Write", []byte{0x30, 0x00}).Return(2, nil)
	mockDevice.On("Read", mock.Anything).Return([]byte{0x1c, 0x02, 'P', '<', 'G', 'B', 'R', 'J', 'O', 'H', 'N', 'S', 'O', 'N', '<', '<', 'J', 'O', 'H', 'N', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<', '<',
*/
