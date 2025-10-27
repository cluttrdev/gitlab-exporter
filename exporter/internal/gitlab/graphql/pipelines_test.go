package graphql

import (
	"testing"
	"time"
)

func TestConvertPipeline(t *testing.T) {
	now := time.Now()

	tests := []struct {
		name    string
		input   PipelineFields
		wantErr bool
	}{
		{
			name: "valid pipeline with core fields",
			input: PipelineFields{
				PipelineReferenceFields: PipelineReferenceFields{
					Id:  "gid://gitlab/Ci::Pipeline/123",
					Iid: "456",
				},
				Project: ProjectReferenceFields{
					Id:       "gid://gitlab/Project/789",
					FullPath: "group/project",
				},
				PipelineFieldsCore: PipelineFieldsCore{
					Name:   ptr("test-pipeline"),
					Ref:    ptr("main"),
					Status: PipelineStatusEnumSuccess,

					CreatedAt:  now,
					UpdatedAt:  now,
					StartedAt:  &now,
					FinishedAt: &now,

					Duration:       ptr(120),
					QueuedDuration: ptr(10.5),
					Coverage:       ptr(85.5),
				},
			},
			wantErr: false,
		},
		{
			name: "pipeline with relations",
			input: PipelineFields{
				PipelineReferenceFields: PipelineReferenceFields{
					Id:  "gid://gitlab/Ci::Pipeline/123",
					Iid: "456",
				},
				Project: ProjectReferenceFields{
					Id:       "gid://gitlab/Project/789",
					FullPath: "group/project",
				},
				PipelineFieldsCore: PipelineFieldsCore{
					Status:    PipelineStatusEnumSuccess,
					CreatedAt: now,
					UpdatedAt: now,
				},
				PipelineFieldsRelations: PipelineFieldsRelations{
					Child: true,
					Upstream: &PipelineFieldsRelationsUpstreamPipeline{
						PipelineReferenceFields: PipelineReferenceFields{
							Id:  "gid://gitlab/Ci::Pipeline/100",
							Iid: "200",
						},
						Project: &PipelineFieldsRelationsUpstreamPipelineProject{
							ProjectReferenceFields: ProjectReferenceFields{
								Id:       "gid://gitlab/Project/300",
								FullPath: "group/upstream-project",
							},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid pipeline id",
			input: PipelineFields{
				PipelineReferenceFields: PipelineReferenceFields{
					Id:  "invalid-id",
					Iid: "456",
				},
				Project: ProjectReferenceFields{
					Id:       "gid://gitlab/Project/789",
					FullPath: "group/project",
				},
				PipelineFieldsCore: PipelineFieldsCore{
					Status:    PipelineStatusEnumSuccess,
					CreatedAt: now,
					UpdatedAt: now,
				},
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ConvertPipeline(tt.input)

			if (err != nil) != tt.wantErr {
				t.Errorf("ConvertPipeline() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err != nil {
				return // Expected error, test passed
			}

			// Verify basic fields
			if result.Id == 0 {
				t.Error("expected non-zero pipeline ID")
			}
			if result.Iid == 0 {
				t.Error("expected non-zero pipeline IID")
			}
			if result.Project.Id == 0 {
				t.Error("expected non-zero project ID")
			}

			// Verify relations if present
			if tt.input.PipelineFieldsRelations.Child {
				if !result.Child {
					t.Error("expected child to be true")
				}
			}
			if tt.input.PipelineFieldsRelations.Upstream != nil {
				if result.UpstreamPipeline == nil {
					t.Error("expected upstream pipeline to be set")
				} else if result.UpstreamPipeline.Id == 0 {
					t.Error("expected non-zero upstream pipeline ID")
				}
			}
		})
	}
}

func TestGetPipelinesOptions(t *testing.T) {
	now := time.Now()

	opts := GetPipelinesOptions{
		UpdatedAfter:  &now,
		UpdatedBefore: &now,
	}

	if opts.UpdatedAfter == nil {
		t.Error("expected UpdatedAfter to be set")
	}
	if opts.UpdatedBefore == nil {
		t.Error("expected UpdatedBefore to be set")
	}
}

// Test that split queries maintain data integrity
func TestPipelineFieldsSplitting(t *testing.T) {
	// This test verifies that core and relations can be combined correctly
	coreFields := PipelineFields{
		PipelineReferenceFields: PipelineReferenceFields{
			Id:  "gid://gitlab/Ci::Pipeline/123",
			Iid: "456",
		},
		Project: ProjectReferenceFields{
			Id:       "gid://gitlab/Project/789",
			FullPath: "group/project",
		},
		PipelineFieldsCore: PipelineFieldsCore{
			Name:   ptr("test-pipeline"),
			Status: PipelineStatusEnumSuccess,

			CreatedAt: time.Now(),
			UpdatedAt: time.Now(),
		},
	}

	relationsFields := PipelineFieldsRelations{
		Child: true,
		Upstream: &PipelineFieldsRelationsUpstreamPipeline{
			PipelineReferenceFields: PipelineReferenceFields{
				Id:  "gid://gitlab/Ci::Pipeline/100",
				Iid: "200",
			},
			Project: &PipelineFieldsRelationsUpstreamPipelineProject{
				ProjectReferenceFields: ProjectReferenceFields{
					Id:       "gid://gitlab/Project/300",
					FullPath: "group/upstream",
				},
			},
		},
	}

	// Simulate merging (like the implementation does)
	combined := coreFields
	combined.PipelineFieldsRelations = relationsFields

	// Convert and verify
	result, err := ConvertPipeline(combined)
	if err != nil {
		t.Fatalf("ConvertPipeline() error = %v", err)
	}

	// Verify core fields are present
	if result.Name != "test-pipeline" {
		t.Errorf("expected name 'test-pipeline', got %v", result.Name)
	}

	// Verify relations are present
	if !result.Child {
		t.Error("expected child to be true")
	}
	if result.UpstreamPipeline == nil {
		t.Error("expected upstream pipeline to be set")
	}
}
