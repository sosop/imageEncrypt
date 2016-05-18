package imageEncrypt

import "testing"

func TestAssembing(t *testing.T) {
	m := NewMetaByRedis("127.0.0.1:6379", "test")
	s := NewFileStorage("test-asserts/")
	a := NewFileSystemAssembe(s, m)
	// _, _, err := a.assembing("test1")
	data, err := a.assebingBase64("test1")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(data)
}
