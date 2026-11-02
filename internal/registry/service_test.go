package registry

import (
	"context"
	"testing"
)

func TestServiceGetAll(t *testing.T) {
	tests := []struct{
		name							string
		expectedResponse  *Registry
		expectedError     error
	} {
		{
			name: "get_all",
			expectedResponse: &Registry{
				Lookup: map[string]int{
					"users": 1,
					"educations": 2,
				},
			},
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService(
				NewMockRegistry(),
				)

			response, err := service.GetAll(t.Context())
			if test.expectedError != err {
				t.Fatalf("unexpected err, %v got %v", test.expectedError, err)
			}

			if test.expectedResponse == nil && response != nil {
				t.Fatalf("unexpected response, %v got %v", test.expectedResponse, response)
			}

			if len(test.expectedResponse.Lookup) != len(response.Lookup) {
				t.Fatalf("unexpected response len, %v got %v", len(test.expectedResponse.Lookup), len(response.Lookup))
			}

			for key, val := range response.Lookup {
				row, ok := test.expectedResponse.Lookup[key]
				if !ok {
					t.Fatalf("unexpected key, %s does not exist", key)
				}

				if row != val {
					t.Fatalf("unexpected value, %v got %v (%s)", row, val, key)
				}
			}
		})
	}
}

func TestServiceCreate(t *testing.T) {
	tests := []struct{
		name							string
		tableName					string
		expectedResponse  *int
		expectedError     error
	} {
		{
			name: "create",
			tableName: "accounts",
			expectedResponse: new(3),
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService(
				NewMockRegistry(),
				)

			response, err := service.Create(t.Context(), test.tableName)
			if test.expectedError != err {
				t.Fatalf("unexpected err, %v got %v", test.expectedError, err)
			}

			if *test.expectedResponse != *response {
				t.Fatalf("unexpected response, %d got %d", *test.expectedResponse, *response)
			}
		})
	}
}

func TestServiceUpdate(t *testing.T) {
	tests := []struct{
		name							string
		id								int
		tableName					string
		expectedError     error
	} {
		{
			name: "update",
			id: 1,
			tableName: "accounts",
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService(
				NewMockRegistry(),
				)

			err := service.Update(t.Context(), test.id, test.tableName)
			if test.expectedError != err {
				t.Fatalf("unexpected err, %v got %v", test.expectedError, err)
			}
		})
	}
}

func TestServiceRemove(t *testing.T) {
	tests := []struct{
		name							string
		id								int
		tableName					string
		expectedError     error
	} {
		{
			name: "remove",
			id: 1,
			expectedError: nil,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			service := NewService(
				NewMockRegistry(),
				)

			err := service.Remove(t.Context(), test.id)
			if test.expectedError != err {
				t.Fatalf("unexpected err, %v got %v", test.expectedError, err)
			}
		})
	}
}

type mockRegistry struct {
	table map[int]string
}

var _ repository = (*mockRegistry)(nil)

func NewMockRegistry() *mockRegistry {
	return &mockRegistry{
		table: map[int]string{
			1: "users",
			2: "educations",
		},
	}
}

func (r *mockRegistry) GetAll(ctx context.Context) (*Registry, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}

	var response = make(map[string]int)

	for key, val := range r.table {
		response[val] = key
	}

	return &Registry{
		Lookup: response,
	}, nil
}

func (r *mockRegistry) Create(ctx context.Context, name string) (*int, error) {
	if ctx.Err() != nil {
		return nil, ctx.Err()
	}
	id := len(r.table) + 1

	r.table[id] = name

	return &id, nil
}

func (r *mockRegistry) Update(ctx context.Context, id int, name string) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	if _, ok := r.table[id]; ok {
		r.table[id] = name
	}

	return nil
}

func (r *mockRegistry) Remove(ctx context.Context, id int) error {
	if ctx.Err() != nil {
		return ctx.Err()
	}

	delete(r.table, id)

	return nil
}
