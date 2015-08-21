package namesys

import (
	"testing"

	key "github.com/heems/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-ipfs/blocks/key"
	path "github.com/heems/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-ipfs/path"
	mockrouting "github.com/heems/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-ipfs/routing/mock"
	u "github.com/heems/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-ipfs/util"
	testutil "github.com/heems/go-ipfs/Godeps/_workspace/src/github.com/ipfs/go-ipfs/util/testutil"
	context "github.com/heems/go-ipfs/Godeps/_workspace/src/golang.org/x/net/context"
)

func TestRoutingResolve(t *testing.T) {
	d := mockrouting.NewServer().Client(testutil.RandIdentityOrFatal(t))

	resolver := NewRoutingResolver(d)
	publisher := NewRoutingPublisher(d)

	privk, pubk, err := testutil.RandTestKeyPair(512)
	if err != nil {
		t.Fatal(err)
	}

	h := path.FromString("/ipfs/QmZULkCELmmk5XNfCgTnCyFgAVxBRBXyDHGGMVoLFLiXEN")
	err = publisher.Publish(context.Background(), privk, h)
	if err != nil {
		t.Fatal(err)
	}

	pubkb, err := pubk.Bytes()
	if err != nil {
		t.Fatal(err)
	}

	pkhash := u.Hash(pubkb)
	res, err := resolver.Resolve(context.Background(), key.Key(pkhash).Pretty())
	if err != nil {
		t.Fatal(err)
	}

	if res != h {
		t.Fatal("Got back incorrect value.")
	}
}
