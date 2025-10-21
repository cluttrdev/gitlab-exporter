package graphql

import (
	"testing"
	"time"
)

func TestConvertRunner(t *testing.T) {
	now := time.Now()
	shortSha := "12345678"

	tests := []struct {
		name    string
		input   RunnerFields
		wantErr bool
	}{
		{
			name: "valid runner with all fields",
			input: RunnerFields{
				RunnerReferenceFields: RunnerReferenceFields{
					Id:       "gid://gitlab/Ci::Runner/123",
					ShortSha: &shortSha,
				},
				RunnerFieldsCore: RunnerFieldsCore{
					Description: ptr("Production Runner"),
					RunnerType:  CiRunnerTypeInstanceType,
					TagList:     []string{"docker", "linux"},
					Status:      CiRunnerStatusOnline,
					Active:      true,
					Locked:      ptr(false),
					Paused:      false,
					AccessLevel: CiRunnerAccessLevelNotProtected,
					RunUntagged: true,
					ContactedAt: &now,
					CreatedAt:   &now,
					CreatedBy: &RunnerFieldsCoreCreatedByUserCore{
						UserReferenceFieldsUserCore: UserReferenceFieldsUserCore{
							Id:       "gid://gitlab/User/456",
							Username: "admin",
							Name:     "Admin User",
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "valid runner without optional fields",
			input: RunnerFields{
				RunnerReferenceFields: RunnerReferenceFields{
					Id:       "gid://gitlab/Ci::Runner/789",
					ShortSha: nil,
				},
				RunnerFieldsCore: RunnerFieldsCore{
					Description: nil,
					RunnerType:  CiRunnerTypeProjectType,
					TagList:     []string{},
					Status:      CiRunnerStatusOffline,
					Active:      false,
					Locked:      nil,
					Paused:      true,
					AccessLevel: CiRunnerAccessLevelRefProtected,
					RunUntagged: false,
					ContactedAt: nil,
					CreatedAt:   nil,
					CreatedBy:   nil,
				},
			},
			wantErr: false,
		},
		{
			name: "runner with group type",
			input: RunnerFields{
				RunnerReferenceFields: RunnerReferenceFields{
					Id:       "gid://gitlab/Ci::Runner/999",
					ShortSha: &shortSha,
				},
				RunnerFieldsCore: RunnerFieldsCore{
					Description: ptr("Group Runner"),
					RunnerType:  CiRunnerTypeGroupType,
					TagList:     []string{"kubernetes"},
					Status:      CiRunnerStatusStale,
					Active:      true,
					Locked:      ptr(true),
					Paused:      false,
					AccessLevel: CiRunnerAccessLevelNotProtected,
					RunUntagged: false,
					ContactedAt: &now,
					CreatedAt:   &now,
					CreatedBy:   nil,
				},
			},
			wantErr: false,
		},
		{
			name: "runner with never contacted status",
			input: RunnerFields{
				RunnerReferenceFields: RunnerReferenceFields{
					Id:       "gid://gitlab/Ci::Runner/111",
					ShortSha: &shortSha,
				},
				RunnerFieldsCore: RunnerFieldsCore{
					Description: ptr("New Runner"),
					RunnerType:  CiRunnerTypeInstanceType,
					TagList:     []string{},
					Status:      CiRunnerStatusNeverContacted,
					Active:      true,
					Locked:      ptr(false),
					Paused:      false,
					AccessLevel: CiRunnerAccessLevelNotProtected,
					RunUntagged: true,
					ContactedAt: nil,
					CreatedAt:   &now,
					CreatedBy:   nil,
				},
			},
			wantErr: false,
		},
		{
			name: "invalid runner id",
			input: RunnerFields{
				RunnerReferenceFields: RunnerReferenceFields{
					Id:       "invalid-id",
					ShortSha: nil,
				},
				RunnerFieldsCore: RunnerFieldsCore{
					RunnerType:  CiRunnerTypeInstanceType,
					Status:      CiRunnerStatusOnline,
					Active:      true,
					AccessLevel: CiRunnerAccessLevelNotProtected,
					RunUntagged: false,
				},
			},
			wantErr: true,
		},
		{
			name: "invalid user id in created by",
			input: RunnerFields{
				RunnerReferenceFields: RunnerReferenceFields{
					Id:       "gid://gitlab/Ci::Runner/123",
					ShortSha: &shortSha,
				},
				RunnerFieldsCore: RunnerFieldsCore{
					RunnerType:  CiRunnerTypeInstanceType,
					Status:      CiRunnerStatusOnline,
					Active:      true,
					AccessLevel: CiRunnerAccessLevelNotProtected,
					RunUntagged: false,
					CreatedBy: &RunnerFieldsCoreCreatedByUserCore{
						UserReferenceFieldsUserCore: UserReferenceFieldsUserCore{
							Id:       "invalid-user-id",
							Username: "admin",
							Name:     "Admin User",
						},
					},
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertRunner(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertRunner() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Expected error, test passed
			}

			// Verify basic fields
			if result.Id == 0 {
				t.Error("expected non-zero runner ID")
			}

			// Verify ShortSha if present
			if tt.input.ShortSha != nil {
				if result.ShortSha != *tt.input.ShortSha {
					t.Errorf("expected ShortSha %v, got %v", *tt.input.ShortSha, result.ShortSha)
				}
			}

			// Verify Description if present
			if tt.input.Description != nil {
				if result.Description != *tt.input.Description {
					t.Errorf("expected Description %v, got %v", *tt.input.Description, result.Description)
				}
			}

			// Verify RunnerType conversion
			if result.RunnerType == "" {
				t.Error("expected non-empty runner type")
			}

			// Verify Status conversion
			if result.Status == "" {
				t.Error("expected non-empty runner status")
			}

			// Verify CreatedBy if present
			if tt.input.CreatedBy != nil {
				if result.CreatedBy.Id == 0 {
					t.Error("expected non-zero user ID")
				}
				if result.CreatedBy.Username != tt.input.CreatedBy.Username {
					t.Errorf("expected username %v, got %v", tt.input.CreatedBy.Username, result.CreatedBy.Username)
				}
				if result.CreatedBy.Name != tt.input.CreatedBy.Name {
					t.Errorf("expected name %v, got %v", tt.input.CreatedBy.Name, result.CreatedBy.Name)
				}
			}

			// Verify boolean fields
			if result.Active != tt.input.Active {
				t.Errorf("expected Active %v, got %v", tt.input.Active, result.Active)
			}
			if result.Paused != tt.input.Paused {
				t.Errorf("expected Paused %v, got %v", tt.input.Paused, result.Paused)
			}
			if result.RunUntagged != tt.input.RunUntagged {
				t.Errorf("expected RunUntagged %v, got %v", tt.input.RunUntagged, result.RunUntagged)
			}

			// Verify TagList
			if len(result.TagList) != len(tt.input.TagList) {
				t.Errorf("expected %d tags, got %d", len(tt.input.TagList), len(result.TagList))
			}
		})
	}
}

func TestConvertRunnerType(t *testing.T) {
	tests := []struct {
		name     string
		input    CiRunnerType
		expected string
	}{
		{"instance type", CiRunnerTypeInstanceType, "INSTANCE"},
		{"group type", CiRunnerTypeGroupType, "GROUP"},
		{"project type", CiRunnerTypeProjectType, "PROJECT"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertRunnerType(tt.input)
			if string(result) != tt.expected {
				t.Errorf("convertRunnerType() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertRunnerStatus(t *testing.T) {
	tests := []struct {
		name     string
		input    CiRunnerStatus
		expected string
	}{
		{"online", CiRunnerStatusOnline, "ONLINE"},
		{"offline", CiRunnerStatusOffline, "OFFLINE"},
		{"stale", CiRunnerStatusStale, "STALE"},
		{"never contacted", CiRunnerStatusNeverContacted, "NEVER_CONTACTED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertRunnerStatus(tt.input)
			if string(result) != tt.expected {
				t.Errorf("convertRunnerStatus() = %v, want %v", result, tt.expected)
			}
		})
	}
}

func TestConvertRunnerAccessLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    CiRunnerAccessLevel
		expected string
	}{
		{"not protected", CiRunnerAccessLevelNotProtected, "NOT_PROTECTED"},
		{"ref protected", CiRunnerAccessLevelRefProtected, "REF_PROTECTED"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := convertRunnerAccessLevel(tt.input)
			if string(result) != tt.expected {
				t.Errorf("convertRunnerAccessLevel() = %v, want %v", result, tt.expected)
			}
		})
	}
}
