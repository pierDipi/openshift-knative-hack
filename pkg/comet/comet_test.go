package comet

import (
	"testing"
)

func TestToComet(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		in      Guess
		want    *Comet
		wantErr bool
	}{
		{
			in: Guess{
				Image: "serverless-operator-135/serverless-kn-operator",
			},
			want: &Comet{
				To: "openshift-serverless-1/serverless-rhel8-operator",
			},
			wantErr: false,
		},
		{
			in: Guess{
				Image: "serverless-operator-135/kn-serving-activator",
			},
			want: &Comet{
				To: "openshift-serverless-1/serving-activator-rhel8",
			},
			wantErr: false,
		},
		{
			in: Guess{
				Image: "kn-serving-activator",
			},
			want: &Comet{
				To: "openshift-serverless-1/serving-activator-rhel8",
			},
			wantErr: false,
		},
		{
			in: Guess{
				Image: "serverless-operator-135/kn-eventing-channel-controller",
			},
			want: &Comet{
				To: "openshift-serverless-1/eventing-in-memory-channel-controller-rhel8",
			},
			wantErr: false,
		},
		{
			in: Guess{
				Image: "serverless-operator-136/kn-eventing-channel-controller",
			},
			want: &Comet{
				To: "openshift-serverless-1/eventing-in-memory-channel-controller-rhel8",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.in.Image, func(t *testing.T) {
			tt := tt

			if tt.in.RHELVersion == "" {
				tt.in.RHELVersion = RHEL8
			}
			if tt.in.FilePath == "" {
				tt.in.FilePath = "comet.yaml"
			}
			got, err := GuessComet(tt.in)
			if (err != nil) != tt.wantErr {
				t.Errorf("GuessComet() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got.To != tt.want.To {
				t.Errorf("GuessComet() got = %q, want %q", got.To, tt.want.To)
			}
		})
	}
}
