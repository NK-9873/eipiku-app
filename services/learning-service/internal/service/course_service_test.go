package service

import (
	"context"
	"reflect"
	"testing"

	"eipiku/learning-service/internal/model"
)

// mockCourseRepository は repository.CourseRepository のモック
type mockCourseRepository struct {
	courses []model.Course
	err     error
}

func (m *mockCourseRepository) FindAll(ctx context.Context) ([]model.Course, error) {
	return m.courses, m.err
}

func TestCourseService_ListCourses(t *testing.T) {
	desc := "身近な動物の名前を覚えよう"
	tests := []struct {
		name string
		repo []model.Course
		want []model.Course
	}{
		{
			name: "repository のコースをそのまま返す",
			repo: []model.Course{{ID: 1, Name: "動物", Description: &desc}, {ID: 2, Name: "果物"}},
			want: []model.Course{{ID: 1, Name: "動物", Description: &desc}, {ID: 2, Name: "果物"}},
		},
		{
			name: "0件なら nil ではなく空のスライスを返す",
			repo: nil,
			want: []model.Course{},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := NewCourseService(&mockCourseRepository{courses: tt.repo})

			got, err := s.ListCourses(context.Background())
			if err != nil {
				t.Fatalf("エラーが返った: %v", err)
			}
			if got == nil || !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got %#v, want %#v", got, tt.want)
			}
		})
	}
}
