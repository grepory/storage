package simple

import "github.com/grepory/storage/apis/meta"

type Simple struct {
	meta.TypeMeta `json:",inline"`

	// Standard object metadata.
	// +optional
	*meta.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Field string `protobuf:"bytes,2,opt,name=field"`
}
