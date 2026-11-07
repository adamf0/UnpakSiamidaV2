package infrastructure

import (
	"fmt"
	"strconv"
	"strings"
	"sync"
	"time"

	"gorm.io/gorm"
)

type MasterIkuItem struct {
	ID           uint64  `gorm:"primaryKey;column:id" json:"id"`
	IdKonteks    uint64  `gorm:"column:id_konteks" json:"id_konteks"`
	Iku          string  `gorm:"column:iku" json:"iku"`
	Kategori     string  `gorm:"column:kategori" json:"kategori"`
	IdSub        *uint64 `gorm:"column:id_sub" json:"id_sub"`
	FullKategori string  `gorm:"column:full_kategori" json:"full_kategori"`
	Satuan       string  `gorm:"column:satuan" json:"satuan"`
	JenisIku     string  `gorm:"column:jenis_iku" json:"jenis_iku"`
	Target       string  `gorm:"column:target" json:"target"`
	FormulaText  string  `gorm:"column:formula_text" json:"formula_text"`
	SumberData   string  `gorm:"column:sumber_data" json:"sumber_data"`
}

var (
	ikuCacheMu   sync.RWMutex
	ikuCacheMap  map[uint64]MasterIkuItem
	ikuCacheList []MasterIkuItem
	ikuCacheTime time.Time
)

func InvalidateIkuCache() {
	ikuCacheMu.Lock()
	defer ikuCacheMu.Unlock()
	ikuCacheMap = nil
	ikuCacheList = nil
	ikuCacheTime = time.Time{}
}

func GetMasterIkuList(dbIku *gorm.DB) []MasterIkuItem {
	if dbIku == nil {
		return []MasterIkuItem{}
	}
	ikuCacheMu.RLock()
	if time.Since(ikuCacheTime) < 1*time.Minute && ikuCacheList != nil {
		defer ikuCacheMu.RUnlock()
		return ikuCacheList
	}
	ikuCacheMu.RUnlock()

	ikuCacheMu.Lock()
	defer ikuCacheMu.Unlock()
	if time.Since(ikuCacheTime) < 1*time.Minute && ikuCacheList != nil {
		return ikuCacheList
	}

	var list []MasterIkuItem
	// Coba tabel master_indikator dulu (tabel native di unpak_iku)
	if err := dbIku.Table("master_indikator").Order("id asc").Find(&list).Error; err != nil {
		_ = dbIku.Table("master_iku").Order("id asc").Find(&list).Error
	}

	newMap := make(map[uint64]MasterIkuItem, len(list))
	for _, item := range list {
		newMap[item.ID] = item
	}
	ikuCacheMap = newMap
	ikuCacheList = list
	ikuCacheTime = time.Now()
	return list
}

func GetMasterIkuMap(dbIku *gorm.DB) map[uint64]MasterIkuItem {
	if dbIku == nil {
		return map[uint64]MasterIkuItem{}
	}
	ikuCacheMu.RLock()
	if time.Since(ikuCacheTime) < 1*time.Minute && ikuCacheMap != nil {
		defer ikuCacheMu.RUnlock()
		return ikuCacheMap
	}
	ikuCacheMu.RUnlock()

	_ = GetMasterIkuList(dbIku)
	ikuCacheMu.RLock()
	defer ikuCacheMu.RUnlock()
	if ikuCacheMap == nil {
		return map[uint64]MasterIkuItem{}
	}
	return ikuCacheMap
}

func FormatIkuList(rawIds string, ikuMap map[uint64]MasterIkuItem) string {
	rawIds = strings.TrimSpace(rawIds)
	if rawIds == "" {
		return ""
	}
	// Jika sudah dalam format "id::iku::jenis", kembalikan apa adanya
	if strings.Contains(rawIds, "::") {
		return rawIds
	}

	parts := strings.Split(rawIds, ",")
	var formatted []string
	for _, p := range parts {
		p = strings.TrimSpace(p)
		if p == "" {
			continue
		}
		if id, err := strconv.ParseUint(p, 10, 64); err == nil {
			if iku, ok := ikuMap[id]; ok {
				formatted = append(formatted, fmt.Sprintf("%d::%s::%s", iku.ID, iku.Iku, iku.JenisIku))
			} else {
				formatted = append(formatted, fmt.Sprintf("%d::IKU %d::", id, id))
			}
		}
	}
	return strings.Join(formatted, "||")
}

func FormatIkuListPtr(rawIds *string, ikuMap map[uint64]MasterIkuItem) *string {
	if rawIds == nil {
		return nil
	}
	res := FormatIkuList(*rawIds, ikuMap)
	return &res
}

