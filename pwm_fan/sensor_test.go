package pwm_fan

import (
	"testing"
)

func TestGetDesiredSpeed(t *testing.T) {
	tempTable := map[float64]float64{
		0:   0,
		30:  50,
		40:  60,
		50:  70,
		60:  80,
		100: 100,
	}

	temps := []float64{100, 60, 50, 40, 30}

	tests := []struct {
		name        string
		currentTemp float64
		want        float64
		wantErr     bool
	}{
		{
			name:        "Test 1",
			currentTemp: 25,
			want:        0,
			wantErr:     true,
		},
		{
			name:        "Test 2",
			currentTemp: 35,
			want:        50,
			wantErr:     false,
		},
		{
			name:        "Test 3",
			currentTemp: 55,
			want:        70,
			wantErr:     false,
		},
		{
			name:        "Test 4",
			currentTemp: 70,
			want:        80,
			wantErr:     false,
		},
		{
			name:        "Test 5",
			currentTemp: 110,
			want:        100,
			wantErr:     false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getDesiredSpeed(tt.currentTemp, temps, tempTable)
			if (err != nil) != tt.wantErr {
				t.Errorf("getDesiredSpeed() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("getDesiredSpeed() = %v, want %v", got, tt.want)
			}
		})
	}
}
