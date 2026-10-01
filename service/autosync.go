package service

import (
	"github.com/lejianwen/rustdesk-api/v2/model"
)

// XC: AutoSyncService mirrors the peer table into a shared address book
// (collection) so administrators always see every device without any
// out-of-tree cron/script. It replaces the legacy /root/rustdesk-ab-sync.py.
//
// Configuration (conf/config.yaml or environment variables):
//
//	rustdesk:
//	  auto-ab-sync: true               # RUSTDESK_API_RUSTDESK_AUTO_AB_SYNC
//	  auto-ab-collection: "全部设备"    # RUSTDESK_API_RUSTDESK_AUTO_AB_COLLECTION
//
// The collection is looked up by name; when it does not exist the sync is a
// no-op (a warning is logged once at startup).
type AutoSyncService struct {
}

// TargetCollection returns the shared collection to sync into, or nil when
// the sync is disabled or the collection cannot be found.
func (as *AutoSyncService) TargetCollection() *model.AddressBookCollection {
	if !Config.Rustdesk.AutoAbSync {
		return nil
	}
	name := Config.Rustdesk.AutoAbCollection
	if name == "" {
		return nil
	}
	c := &model.AddressBookCollection{}
	DB.Where("name = ?", name).First(c)
	if c.Id == 0 {
		return nil
	}
	return c
}

// SyncPeer mirrors one peer into the target collection. Existing entries only
// have machine-reported fields refreshed - aliases, tags, passwords and other
// user customisations are never overwritten. Missing entries are created.
func (as *AutoSyncService) SyncPeer(peer *model.Peer) {
	c := as.TargetCollection()
	if c == nil || peer == nil || peer.Id == "" {
		return
	}
	ab := &model.AddressBook{}
	DB.Where("collection_id = ? and id = ?", c.Id, peer.Id).First(ab)
	if ab.RowId == 0 {
		n := AllService.AddressBookService.FromPeer(peer)
		n.CollectionId = c.Id
		n.UserId = c.UserId
		if err := DB.Create(n).Error; err != nil {
			Logger.Errorf("XC auto-ab-sync: add peer %s to collection %q failed: %v", peer.Id, c.Name, err)
			return
		}
		Logger.Infof("XC auto-ab-sync: added peer %s (%s) to collection %q", peer.Id, peer.Hostname, c.Name)
		return
	}
	DB.Model(ab).Updates(map[string]interface{}{
		"hostname": peer.Hostname,
		"username": peer.Username,
		"platform": AllService.AddressBookService.PlatformFromOs(peer.Os),
	})
}

// SyncAll reconciles every known peer once; used at startup.
func (as *AutoSyncService) SyncAll() {
	if as.TargetCollection() == nil {
		Logger.Warnf("XC auto-ab-sync: enabled but collection %q not found - sync inactive", Config.Rustdesk.AutoAbCollection)
		return
	}
	var peers []*model.Peer
	DB.Find(&peers)
	for _, p := range peers {
		as.SyncPeer(p)
	}
	Logger.Infof("XC auto-ab-sync: reconciliation done (%d peers)", len(peers))
}

// Start kicks off the initial reconciliation in the background.
func (as *AutoSyncService) Start() {
	go as.SyncAll()
}
