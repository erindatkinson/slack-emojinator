package slack

import (
	"context"

	"github.com/browserutils/kooky"
	_ "github.com/browserutils/kooky/browser/all"
)

/*
CookieStoreProcessor

Loops through all found cookie stores and performs the given action on each
*/
func CookieStoreProcessor(ctx context.Context, action func(store kooky.CookieStore)) {
	stores := kooky.FindAllCookieStores(ctx)
	for _, store := range stores {
		action(store)
	}

}
