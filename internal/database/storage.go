// Пакет storage включает в себя функции для работы с БД.
package storage

import (
	"encoding/hex"
	"encoding/json"
	"errors"
	"log"
	"time"

	"github.com/boltdb/bolt"

	MyLogicClient "gophkeeper/internal/logic_client"
	MyLogicServer "gophkeeper/internal/logic_server"
)

// OpenDB для открытия БД.
func OpenDB(path string) *bolt.DB {
	db, err := bolt.Open(path, 0600, nil)
	if err != nil {
		log.Println(err)
	}

	return db
}

// RegisterDB Для записи регистрвционных данных в БД.
func RegisterDB(db *bolt.DB, mlpss MyLogicServer.MailPass, ckslt MyLogicServer.CookieSalt) (error, string) {
	var PassCookieSalt MyLogicServer.CookieSalt = ckslt

	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("MailPassCookieSalt"))
		if err != nil {
			log.Println(err)
			return err
		}

		bmlpss, err := json.Marshal(mlpss)
		if err != nil {
			log.Println(err)
			return err
		}

		bckslt, err := json.Marshal(ckslt)
		if err != nil {
			log.Println(err)
			return err
		}

		err = bucket.ForEach(func(key, value []byte) error {
			var loginPass MyLogicServer.MailPass
			if err := json.Unmarshal(key, &loginPass); err != nil {
				return err
			}
			if loginPass.Mail == mlpss.Mail {
				err = errors.New("User jast exit")
				return err
			}
			return nil
		})
		if err != nil {
			return err
		}

		err = bucket.Put(bmlpss, bckslt)
		if err != nil {
			return err
		}

		return nil
	})
	TimeMark(db)
	if err != nil {
		return err, ""
	}
	return nil, PassCookieSalt.Cookie
}

// GetCookieDB Забарать cookie из БД по паролю и логину.
func GetCookieDB(db *bolt.DB, mlpss MyLogicServer.MailPass) (MyLogicServer.CookieSalt, error) {
	var PassCookieSalt MyLogicServer.CookieSalt

	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("MailPassCookieSalt"))
		if bucket == nil {
			err := errors.New("MailPassCookieSalt Not found")
			return err
		}

		bmlpss, err := json.Marshal(mlpss)
		if err != nil {
			log.Println(err)
			return err
		}

		value := bucket.Get(bmlpss)
		if value != nil {
			if err := json.Unmarshal(value, &PassCookieSalt); err != nil {
				return err
			}
			return nil
		} else {
			err = errors.New("User not exit")
			return err
		}

	})
	if err != nil {
		return PassCookieSalt, err
	}

	return PassCookieSalt, nil
}

// Autorization Функция проверки авторизации пользователя.
func Authorization(db *bolt.DB, cookieValue string) bool {
	Found := false
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("MailPassCookieSalt"))
		if bucket == nil {
			err := errors.New("Bucket nor exist")
			return err
		}

		err := bucket.ForEach(func(key, value []byte) error {
			var cookieSalt MyLogicServer.CookieSalt
			if err := json.Unmarshal(value, &cookieSalt); err != nil {
				return err
			}

			if cookieSalt.Cookie == cookieValue {
				Found = true
			}
			return nil
		})
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return false
	}
	return Found
}

// EditMasterPass Функция изменения мастер пароля в БД.
func EditMasterPass(db *bolt.DB, cookieGet string, newPass string) {
	var CookieSalt MyLogicServer.CookieSalt
	var LoginPass MyLogicServer.MailPass

	log.Println(cookieGet, newPass)

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("MailPassCookieSalt"))
		if bucket == nil {
			err := errors.New("Bucket nor exist")
			return err
		}

		err := bucket.ForEach(func(key, value []byte) error {
			if err := json.Unmarshal(value, &CookieSalt); err != nil {
				return err
			}
			if CookieSalt.Cookie == cookieGet {
				if err := json.Unmarshal(key, &LoginPass); err != nil {
					return err
				}

			}
			return nil
		})
		if err != nil {
			return err
		}

		bmlpss, err := json.Marshal(LoginPass)
		if err != nil {
			log.Println(err)
			return err
		}

		err = bucket.Delete(bmlpss)
		if err != nil {
			return err
		}

		LoginPass.Pass = newPass

		bmlpss, err = json.Marshal(LoginPass)
		if err != nil {
			log.Println(err)
			return err
		}

		bckslt, err := json.Marshal(CookieSalt)
		if err != nil {
			log.Println(err)
			return err
		}

		log.Println("EditMasterPass")
		log.Println(LoginPass.Mail, LoginPass.Pass)
		log.Println(CookieSalt.Cookie, CookieSalt.Salt)
		log.Println("--------------------------------")

		err = bucket.Put(bmlpss, bckslt)
		if err != nil {
			return err
		}
		return nil
	})
	TimeMark(db)
	if err != nil {
		log.Println(err)
	}
}

// AddData Функция добавления данных в БД.
func AddData(db *bolt.DB, cookieGet string, data []byte) int {
	var CookieID MyLogicServer.IDCookie
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("UsersData"))
		if err != nil {
			log.Println(err)
			return err
		}
		bucket = tx.Bucket([]byte("UsersData"))

		var LocalID int = -1
		err = bucket.ForEach(func(key, value []byte) error {
			if err := json.Unmarshal(key, &CookieID); err != nil {
				return err
			}
			if CookieID.Cookie == cookieGet {
				if LocalID < CookieID.ID {
					LocalID = CookieID.ID
				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		CookieID.Cookie = cookieGet
		CookieID.ID = LocalID + 1

		bkey, err := json.Marshal(CookieID)
		if err != nil {
			return err
		}

		err = bucket.Put(bkey, data)
		if err != nil {
			return err
		}

		return nil
	})
	TimeMark(db)
	if err != nil {
		return -1
	}
	return CookieID.ID
}

// UpdateData Функция обновления данных в БД.
func UpdateData(db *bolt.DB, cookieGet string, Data []byte, id int) {
	var CookieID MyLogicServer.IDCookie
	var DataInfo MyLogicClient.MyData
	ChangeData := false
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("UsersData"))
		if err != nil {
			log.Println(err)
			return err
		}

		err = bucket.ForEach(func(key, value []byte) error {
			if err := json.Unmarshal(key, &CookieID); err != nil {
				return err
			}

			if CookieID.Cookie == cookieGet {
				if CookieID.ID == id {
					ChangeData = true
					if err := json.Unmarshal(value, &DataInfo); err != nil {
						return err
					}

				}
			}
			return nil
		})
		if err != nil {
			return err
		}

		if ChangeData {
			CookieID.Cookie = cookieGet
			CookieID.ID = id
			bkey, err := json.Marshal(CookieID)
			if err != nil {
				return err
			}

			DataInfo.BytesData = hex.EncodeToString(Data)

			bvalue, err := json.Marshal(DataInfo)
			if err != nil {
				return err
			}

			err = bucket.Put(bkey, bvalue)
			if err != nil {
				return err
			}
		}
		return nil
	})
	TimeMark(db)
	if err != nil {
		log.Println(err)
	}
}

// GetData Функция взятия данных из БД.
func GetData(db *bolt.DB, cookieGet string, id int) ([]byte, error) {
	var CookieID MyLogicServer.IDCookie
	var DataInfo MyLogicClient.MyData
	err := db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("UsersData"))
		if bucket == nil {
			err := errors.New("Bucket nor exist")
			return err
		}

		err := bucket.ForEach(func(key, value []byte) error {
			if err := json.Unmarshal(key, &CookieID); err != nil {
				log.Println(err)
				return err
			}

			if CookieID.Cookie == cookieGet {
				if CookieID.ID == id {
					if err := json.Unmarshal(value, &DataInfo); err != nil {
						log.Println(err)
						return err
					}
				}
			}
			return nil
		})
		if err != nil {
			log.Println(err)
			return err
		}
		return nil
	})

	bRes, err := json.Marshal(DataInfo)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}
	return bRes, nil
}

// Delete Функция удаления данных из БД.
func Delete(db *bolt.DB, cookieGet string, id int) error {
	CookieID := MyLogicServer.IDCookie{
		ID:     id,
		Cookie: cookieGet,
	}

	err := db.Update(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("UsersData"))
		if bucket == nil {
			err := errors.New("Bucket nor exist")
			return err
		}

		bkey, err := json.Marshal(CookieID)
		if err != nil {
			return err
		}

		err = bucket.Delete(bkey)
		if err != nil {
			return err
		}
		return nil
	})
	TimeMark(db)
	if err != nil {
		return err
	}
	return nil
}

// Тип MyTime необходим для перевода времени в биты, для внесения в БД
type MyTime struct {
	Time int64 `json:"time"`
}

// ТimeМark функция создания и изменения временной метки в БД.
func TimeMark(db *bolt.DB) {
	err := db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("TimeMark"))
		if err != nil {
			return err
		}

		today := MyTime{
			Time: time.Now().Unix(),
		}

		bToday, err := json.Marshal(today)
		if err != nil {
			return nil
		}

		err = bucket.Put([]byte("LastEdit"), bToday)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		log.Println(err)
	}
}

// Closer функция закрытия БД для клиента и сервера.
func Closer(db1 *bolt.DB, db2 *bolt.DB) {
	db1.Close()
	db2.Close()
}

// Sync функция синхронизации БД клиента и сервера.
func Sync(clientdb string, serverdb string) {
	var ClientTime MyTime
	var ServerTime MyTime

	db1 := OpenDB(clientdb)
	db2 := OpenDB(serverdb)
	defer Closer(db1, db2)

	err := db1.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("TimeMark"))
		if bucket == nil {
			err := errors.New("Bucket not found db1 1")
			return err
		}

		value := bucket.Get([]byte("LastEdit"))
		if value != nil {
			if err := json.Unmarshal(value, &ClientTime); err != nil {
				return err
			}
			return nil
		}
		err := errors.New("Value equal 0")
		return err
	})
	if err != nil {
		log.Println(err)
	}

	err = db2.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("TimeMark"))
		if bucket == nil {
			err := errors.New("Bucket not found 1")
			return err
		}

		value := bucket.Get([]byte("LastEdit"))
		if value != nil {
			if err := json.Unmarshal(value, &ServerTime); err != nil {
				return err
			}
			return nil
		}
		err := errors.New("Value equal 0")
		return err
	})
	if err != nil {
		log.Println(err)
	}

	if ClientTime.Time >= ServerTime.Time {
		err = db1.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("MailPassCookieSalt"))
			if bucket == nil {
				err := errors.New("Bucket not found db1 2")
				return err
			}

			err := bucket.ForEach(func(key, value []byte) error {
				err := db2.Update(func(tx *bolt.Tx) error {
					targetBucket, err := tx.CreateBucketIfNotExists([]byte("MailPassCookieSalt"))
					if err != nil {
						return err
					}
					targetBucket = tx.Bucket([]byte("MailPassCookieSalt"))

					err = targetBucket.Put(key, value)
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					log.Println(err)
					return err
				}
				return nil
			})
			if err != nil {
				log.Println(err)
				return err
			}

			return nil
		})
		if err != nil {
			log.Println(err)
		}

		err = db1.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("UsersData"))
			if bucket == nil {
				err := errors.New("Bucket not found db1 3")
				return err
			}

			err := bucket.ForEach(func(key, value []byte) error {
				err := db2.Update(func(tx *bolt.Tx) error {
					targetBucket, err := tx.CreateBucketIfNotExists([]byte("UsersData"))
					if err != nil {
						log.Println(err)
						return err
					}

					err = targetBucket.Put(key, value)
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				})
				if err != nil {
					log.Println(err)
					return err
				}
				return nil
			})
			if err != nil {
				log.Println(err)
				return err
			}

			return nil
		})
		if err != nil {
			log.Println(err)
		}

		err = db1.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("TimeMark"))
			if bucket == nil {
				err := errors.New("Bucket not found db1 3a")
				return err
			}

			err := bucket.ForEach(func(key, value []byte) error {
				err := db2.Update(func(tx *bolt.Tx) error {
					targetBucket, err := tx.CreateBucketIfNotExists([]byte("TimeMark"))
					if err != nil {
						log.Println(err)
						return err
					}

					err = targetBucket.Put(key, value)
					if err != nil {
						log.Println(err)
						return err
					}
					return nil
				})
				if err != nil {
					log.Println(err)
					return err
				}
				return nil
			})
			if err != nil {
				log.Println(err)
				return err
			}

			return nil
		})
		if err != nil {
			log.Println(err)
		}

	} else {
		err = db2.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("MailPassCookieSalt"))
			if bucket == nil {
				err := errors.New("Bucket not found 2")
				return err
			}

			err := bucket.ForEach(func(key, value []byte) error {
				err := db1.Update(func(tx *bolt.Tx) error {
					targetBucket, err := tx.CreateBucketIfNotExists([]byte("MailPassCookieSalt"))

					if err != nil {
						return err
					}
					err = targetBucket.Put(key, value)
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			log.Println(err)
		}

		err = db2.View(func(tx *bolt.Tx) error {
			bucket := tx.Bucket([]byte("UsersData"))
			if bucket == nil {
				err := errors.New("Bucket not found 3")
				return err
			}

			err := bucket.ForEach(func(key, value []byte) error {
				err := db1.Update(func(tx *bolt.Tx) error {
					targetBucket, err := tx.CreateBucketIfNotExists([]byte("UsersData"))

					if err != nil {
						return err
					}
					err = targetBucket.Put(key, value)
					if err != nil {
						return err
					}
					return nil
				})
				if err != nil {
					return err
				}
				return nil
			})
			if err != nil {
				return err
			}

			return nil
		})
		if err != nil {
			log.Println(err)
		}
	}
}
