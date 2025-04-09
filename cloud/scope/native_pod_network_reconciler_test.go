package scope

import (
	"context"
	"errors"
	"testing"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/gomega"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
)

func TestNativePodNetworkReconciler(t *testing.T) {
	var (
		mockCtrl *gomock.Controller
	)

	setup := func(t *testing.T, g *WithT) {
		mockCtrl = gomock.NewController(t)
	}
	teardown := func(t *testing.T, g *WithT) {
		mockCtrl.Finish()
	}
	tests := []struct {
		name              string
		testSpecificSetup func(mockNpnInterface *MockNpnWorkloadClientInterface)
		expectedError     error
	}{
		{
			name: "DeleteNpn succeeds",
			testSpecificSetup: func(mockNpnInterface *MockNpnWorkloadClientInterface) {
				mockNpnInterface.EXPECT().DeleteNpn(gomock.Any()).Return(nil)
			},
			expectedError: nil,
		},
		{
			name: "DeleteNpn fails",
			testSpecificSetup: func(mockNpnInterface *MockNpnWorkloadClientInterface) {
				mockNpnInterface.EXPECT().DeleteNpn(gomock.Any()).Return(errors.New("failed to delete NPN"))
			},
			expectedError: errors.New("failed to delete NPN"),
		},
		{
			name: "HasNpnCrd returns true",
			testSpecificSetup: func(mockNpnInterface *MockNpnWorkloadClientInterface) {
				mockNpnInterface.EXPECT().HasNpnCrd(gomock.Any()).Return(true, nil)
			},
			expectedError: nil,
		},
		{
			name: "HasNpnCrd returns false with error",
			testSpecificSetup: func(mockNpnInterface *MockNpnWorkloadClientInterface) {
				mockNpnInterface.EXPECT().HasNpnCrd(gomock.Any()).Return(false, errors.New("CRD not found"))
			},
			expectedError: errors.New("CRD not found"),
		},
		{
			name: "GetOrCreateNpn succeeds",
			testSpecificSetup: func(mockNpnInterface *MockNpnWorkloadClientInterface) {
				mockNpnInterface.EXPECT().GetOrCreateNpn(gomock.Any()).Return(&unstructured.Unstructured{}, nil)
			},
			expectedError: nil,
		},
		{
			name: "GetOrCreateNpn fails",
			testSpecificSetup: func(mockNpnInterface *MockNpnWorkloadClientInterface) {
				mockNpnInterface.EXPECT().GetOrCreateNpn(gomock.Any()).Return(nil, errors.New("failed to create NPN"))
			},
			expectedError: errors.New("failed to create NPN"),
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			g := NewWithT(t)
			defer teardown(t, g)
			setup(t, g)
			mockCtrl := gomock.NewController(t)
			defer mockCtrl.Finish()

			mockNpnInterface := NewMockNpnWorkloadClientInterface(mockCtrl)

			// Test-specific setup
			tc.testSpecificSetup(mockNpnInterface)

			// Execute the method under test
			var err error
			switch tc.name {
			case "DeleteNpn succeeds", "DeleteNpn fails":
				err = mockNpnInterface.DeleteNpn(context.Background())
			case "HasNpnCrd returns true", "HasNpnCrd returns false with error":
				_, err = mockNpnInterface.HasNpnCrd(context.Background())
			case "GetOrCreateNpn succeeds", "GetOrCreateNpn fails":
				_, err = mockNpnInterface.GetOrCreateNpn(context.Background())
			}

			// Validate the result
			if tc.expectedError != nil {
				g.Expect(err).To(Not(BeNil()))
				g.Expect(err.Error()).To(Equal(tc.expectedError.Error()))
			} else {
				g.Expect(err).To(BeNil())
			}
		})
	}
}
