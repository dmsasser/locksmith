package data

import (
	"bufio"
	//	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"os"
	//	"fmt"
	"testing"
)

func checke(t *testing.T, e error) {
	if e != nil {
		t.Error(e)
	}
}

func assertStringsEquals(t *testing.T, s1, s2 string) {
	if s1 != s2 {
		t.Logf("Expected [%s] but got [%s]", s1, s2)
		t.Fail()
	}
}

func TestSSHPublicKeyParsing(t *testing.T) {
	path := "test-data/public-keys"
	keys, err := ioutil.ReadDir(path)
	checke(t, err)

	for _, key := range keys {
		kp := path + "/" + key.Name()
		t.Run(key.Name(), func(t *testing.T) {
			k := Read(kp)
			if k == nil {
				t.Error("Failed to parse " + kp)
			}
		})
	}
}

func SkipTestSSHPrivateKeyParse(t *testing.T) {
	t.Run("RSA", func(t *testing.T) {
		key := Read("test-data/rsa")
		if key == nil {
			t.Error("Failed to parse RSA")
		}
		rsa := key.(*SSHKey)
		assertStringsEquals(t, "ssh-rsa", rsa.KeyType())
		assertStringsEquals(t, "AAAAB3NzaC1yc2EAAAADAQABAAABAQDEhoo9i/AwdwWx2xFcQjZkQxlNlex1p7pyOn7qitncnc/+bEHSARGoflqMMFgoBMrsKcQUZXt+LpBvlwGbTqATfat5SwKJbQi2EcoRr8j0e1gsG357zv0i/GuemdTctyk2Hdxq+MkuSlSMlswoAPLfGhFBUiBNLIrb5wwK8MNJjpRkqONxtDQHYpeZ7J+PdSVAQYJ6aNxrA5zRd732CHDyMkHIvnmb+vFa7rPYYwLyzborMrTEQXc1IpqNOzkF33AXAmqsjwNabmReRyerVGZ5cyLJEhn0Yjkixa1lt4RcioV8y4OnLXeHOB7DP1HEko3Ox8Tc16r+b2v70+YBc2c5", rsa.PublicKeyString())
	})

	t.Run("DSS", func(t *testing.T) {
		key := Read("test-data/dss")
		if key == nil {
			t.Error("Failed to parse DSS")
		}
		dss := key.(*SSHKey)
		assertStringsEquals(t, "ssh-dss", dss.KeyType())
		assertStringsEquals(t, "AAAAB3NzaC1kc3MAAACBAMZhAjMPsL/oo9RZiD7jfWBOVGoLqwdwtjuTkaKVFmBVBh+c2nMi11zVzRz1JqbXR15QNyaDc2EumZTC2WTyas4uSXTh2F6Ohto+a2QnCN3rjsiBsXHnr6hbBN+Qs8uJ/+ssGDpsWKIpWOL3+Q6QmHQZg+df4XtBlMyehCWr7jCdAAAAFQCrynAE+Z6tGteawaHWa8ReOpYkrQAAAIB3cd1Ls/1ox/gNNMqTbuAvWQIgIda7Uw+OHU55EyeryPR9e2GH6rsHWCwd47cyurOukqF+e5FH/dnj7K/Kt4BFXPeR0YU4KaiAZIEl8I7Kcdazxz3vWgK3sTKRy10ABqEZL9oUazMfX43IaiPeiU6nwgrMHokTwKLkZH+iBwN8JQAAAIEAo+h6Lop9my2BxrHKSmhQfya3rl0N35ZDk/8kExLW1xkpQmzARrCMrw3YNuRCNgrh5Ds7EdyG0HyjWnnSnPBXqCxFfDTtaGeieLquocEK3M5DGckgI4IEa9pvL3fVZ/cHT3YxC369PF/vX9l7TPHF6Au8lnEFEzNyZLQvsfrqxgg=", dss.PublicKeyString())
	})
}

func TestSSHPublicKeyParse(t *testing.T) {
	t.Run("RSA", func(t *testing.T) {
		key := Read("test-data/rsa.pub")
		if key == nil {
			t.Error("Failed to parse RSA")
		}
		rsa := key.(*SSHKey)
		assertStringsEquals(t, "ssh-rsa", rsa.KeyType())
		assertStringsEquals(t, "dewey@FlynnRyder", rsa.Comments[0])
		assertStringsEquals(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQDEhoo9i/AwdwWx2xFcQjZkQxlNlex1p7pyOn7qitncnc/+bEHSARGoflqMMFgoBMrsKcQUZXt+LpBvlwGbTqATfat5SwKJbQi2EcoRr8j0e1gsG357zv0i/GuemdTctyk2Hdxq+MkuSlSMlswoAPLfGhFBUiBNLIrb5wwK8MNJjpRkqONxtDQHYpeZ7J+PdSVAQYJ6aNxrA5zRd732CHDyMkHIvnmb+vFa7rPYYwLyzborMrTEQXc1IpqNOzkF33AXAmqsjwNabmReRyerVGZ5cyLJEhn0Yjkixa1lt4RcioV8y4OnLXeHOB7DP1HEko3Ox8Tc16r+b2v70+YBc2c5\n", rsa.PublicKeyString())
	})

	t.Run("DSS", func(t *testing.T) {
		key := Read("test-data/dss.pub")
		if key == nil {
			t.Error("Failed to parse DSS")
		}
		dss := key.(*SSHKey)
		assertStringsEquals(t, "ssh-dss", dss.KeyType())
		assertStringsEquals(t, "dewey@FlynnRyder", dss.Comments[0])
		assertStringsEquals(t, "ssh-dss AAAAB3NzaC1kc3MAAACBAMZhAjMPsL/oo9RZiD7jfWBOVGoLqwdwtjuTkaKVFmBVBh+c2nMi11zVzRz1JqbXR15QNyaDc2EumZTC2WTyas4uSXTh2F6Ohto+a2QnCN3rjsiBsXHnr6hbBN+Qs8uJ/+ssGDpsWKIpWOL3+Q6QmHQZg+df4XtBlMyehCWr7jCdAAAAFQCrynAE+Z6tGteawaHWa8ReOpYkrQAAAIB3cd1Ls/1ox/gNNMqTbuAvWQIgIda7Uw+OHU55EyeryPR9e2GH6rsHWCwd47cyurOukqF+e5FH/dnj7K/Kt4BFXPeR0YU4KaiAZIEl8I7Kcdazxz3vWgK3sTKRy10ABqEZL9oUazMfX43IaiPeiU6nwgrMHokTwKLkZH+iBwN8JQAAAIEAo+h6Lop9my2BxrHKSmhQfya3rl0N35ZDk/8kExLW1xkpQmzARrCMrw3YNuRCNgrh5Ds7EdyG0HyjWnnSnPBXqCxFfDTtaGeieLquocEK3M5DGckgI4IEa9pvL3fVZ/cHT3YxC369PF/vX9l7TPHF6Au8lnEFEzNyZLQvsfrqxgg=\n", dss.PublicKeyString())
	})
}

func TestReadAuthorizedKeys(t *testing.T) {
	file, err := os.Open("test-data/authorized_keys")
	checke(t, err)
	defer file.Close()

	scanner := bufio.NewScanner(file)

	// Constrained key is the first line
	scanner.Scan()
	text := scanner.Text()

	key := New(text)

	if key == nil {
		t.Errorf("Key is nil")
	}

	sshkey := key.(*SSHKey)

	assertStringsEquals(t, "ssh-rsa", sshkey.KeyType())
	assertStringsEquals(t, "ssh-rsa AAAAB3NzaC1yc2EAAAADAQABAAABAQCpVaCpLlQ8Wf4KgcRwsIvXCvf0Onkp9hZ2Sov5s2ZiIqJne8Rk9kx3CoSHGMpCSCBuGybs8k/8ga7g/l6+bKDc3aDGWw52+7ClBGz4xjL5C9HXub2iKRdxIesDtkQtQMawFobBTi9hiW92SoK1H/AmLhHDxicfidXPaOcNY57PWZDqEmR2PWo0k4oNn0zxQO3UJmfiKNoR6ozJ3JDWGCu2SMh/YobKwNSlge6YsVKO4zpxR3wBbHS9CYL2xE6QMyN1KnJ+ACoeZF8tkXThOAgH5VERoM+KawAHK0Hqpqh8d85jQU7ul9ernFCip2zVAC/hsobORmHGyvGd9aWDXZTB\n", sshkey.PublicKeyString())
	assertStringsEquals(t, "constrained key", sshkey.Comments[0])
	//	assertStringsEquals(t, "command=\"/bin/ps -ef\",no-port-forwarding,no-X11-forwarding,no-pty", sshkey.Constraints)
}

func TestSSHId(t *testing.T) {
	key := Read("test-data/rsa.pub")
	assertStringsEquals(t, "SHA256:mbhMXOdSermDODXkg5fBUQN9yst7W9Fkn9yurscQSOQ", string(key.Id()))
}

func TestSSHJSon(t *testing.T) {
	key := Read("test-data/rsa.pub")

	json, error := key.Json()
	checke(t, error)
	ioutil.WriteFile("temp", json, 666)
	assertStringsEquals(t, `{
  "Type": "SSHKey",
  "Ids": [
    "SHA256:mbhMXOdSermDODXkg5fBUQN9yst7W9Fkn9yurscQSOQ"
  ],
  "Names": [],
  "Deprecated": false,
  "Replacement": "",
  "PublicKey": {
    "N": 24809051920830052890525181663834441562026965135039612171555457049570978341010824702982333278059308143991930998126783941620692557632728744577159151723229554430031053565344135388348118543421038378400575227281198627317117014484306016930895724285350616040857348021361132827891442508498492344373669032945937823849182795792657014842937263625871476146281616750298221384794570364753932252429117402688216047768632943893167376694869582348687268093872504134257344361532896936336925442773163484682792603660177457890027949342950405164971002422265080740838135765734774882286316380152623151063236328138603783949035635609218432853817,
    "E": 65537
  },
  "Comments": [
    "dewey@FlynnRyder"
  ]
}`, string(json))

	newkey := ReadJson(json)

	if key.Id() != newkey.Id() {
		t.Errorf("Failed to load JSON")
	}
}