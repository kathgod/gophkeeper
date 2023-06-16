package storage_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/boltdb/bolt"

	MyStorage "gophkeeper/internal/database"
	MyLogicServer "gophkeeper/internal/logic_server"
)

const (
	projectpath = "../../"
	dbPath      = projectpath + "/internal/database/testdb.db"
)

func TestOpenDB(t *testing.T) {
	dbPath := projectpath + "/internal/database/testdb.db"

	result := MyStorage.OpenDB(dbPath)
	defer os.Remove(dbPath)
	defer result.Close()

	if result == nil {
		t.Error("Expected non-nil database pointer, got nil")
	}

	expectedPath := projectpath + "/internal/database/testdb.db"
	if result.Path() != expectedPath {
		t.Errorf("Expected database path %s, got %s", expectedPath, result.Path())
	}
}

func TestRegisterDB(t *testing.T) {
	dbPath := projectpath + "/internal/database/testdb.db"
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open test database: %s", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	mlpss := MyLogicServer.MailPass{
		Mail: "test@example.com",
		Pass: "testpass",
	}
	ckslt := MyLogicServer.CookieSalt{
		Cookie: "testcookie",
		Salt:   "testsalt",
	}

	err, cookie := MyStorage.RegisterDB(db, mlpss, ckslt)
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	expectedCookie := ckslt.Cookie
	if cookie != expectedCookie {
		t.Errorf("Expected cookie %s, got %s", expectedCookie, cookie)
	}
}

func TestGetCookieDB(t *testing.T) {
	dbPath := projectpath + "/internal/database/testdb.db"
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open test database: %s", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	mlpss := MyLogicServer.MailPass{
		Mail: "test@example.com",
		Pass: "testpass",
	}
	ckslt := MyLogicServer.CookieSalt{
		Cookie: "testcookie",
		Salt:   "testsalt",
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("MailPassCookieSalt"))
		if err != nil {
			return err
		}

		bmlpss, err := json.Marshal(mlpss)
		if err != nil {
			return err
		}

		bckslt, err := json.Marshal(ckslt)
		if err != nil {
			return err
		}

		err = bucket.Put(bmlpss, bckslt)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed to initialize test data: %s", err)
	}

	result, err := MyStorage.GetCookieDB(db, mlpss)
	if err != nil {
		t.Fatalf("Expected no error, got: %s", err)
	}

	expectedResult := ckslt
	if result != expectedResult {
		t.Errorf("Expected result %+v, got %+v", expectedResult, result)
	}
}

func TestAuthorization(t *testing.T) {
	dbPath := projectpath + "/internal/database/testdb.db"
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open test database: %s", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	mlpss := MyLogicServer.MailPass{
		Mail: "test@example.com",
		Pass: "testpass",
	}
	ckslt := MyLogicServer.CookieSalt{
		Cookie: "testcookie",
		Salt:   "testsalt",
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("MailPassCookieSalt"))
		if err != nil {
			return err
		}

		bmlpss, err := json.Marshal(mlpss)
		if err != nil {
			return err
		}

		bckslt, err := json.Marshal(ckslt)
		if err != nil {
			return err
		}

		err = bucket.Put(bmlpss, bckslt)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed to initialize test data: %s", err)
	}

	result := MyStorage.Authorization(db, ckslt.Cookie)
	if result != true {
		t.Error("Expected true, got false")
	}

	result = MyStorage.Authorization(db, "invalidcookie")
	if result != false {
		t.Error("Expected false, got true")
	}
}

func TestEditMasterPass(t *testing.T) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open test database: %s", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	mlpss := MyLogicServer.MailPass{
		Mail: "test@example.com",
		Pass: "testpass",
	}
	ckslt := MyLogicServer.CookieSalt{
		Cookie: "testcookie",
		Salt:   "testsalt",
	}

	err = db.Update(func(tx *bolt.Tx) error {
		bucket, err := tx.CreateBucketIfNotExists([]byte("MailPassCookieSalt"))
		if err != nil {
			return err
		}

		bmlpss, err := json.Marshal(mlpss)
		if err != nil {
			return err
		}

		bckslt, err := json.Marshal(ckslt)
		if err != nil {
			return err
		}

		err = bucket.Put(bmlpss, bckslt)
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed to initialize test data: %s", err)
	}

	MyStorage.EditMasterPass(db, ckslt.Cookie, "newpass")

	var updatedLoginPass MyLogicServer.MailPass
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("MailPassCookieSalt"))
		if bucket == nil {
			err := errors.New("Bucket not exist")
			return err
		}

		err := bucket.ForEach(func(key, value []byte) error {
			if err := json.Unmarshal(key, &updatedLoginPass); err != nil {
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
		t.Fatalf("Failed to retrieve updated data: %s", err)
	}

	expectedPass := "newpass"
	if updatedLoginPass.Pass != expectedPass {
		t.Errorf("Expected updated password %s, got %s", expectedPass, updatedLoginPass.Pass)
	}
}

func TestAddData(t *testing.T) {
	db, err := bolt.Open(dbPath, 0600, nil)
	if err != nil {
		t.Fatalf("Failed to open test database: %s", err)
	}
	defer os.Remove(dbPath)
	defer db.Close()

	cookie := "testcookie"
	data := []byte("testdata")

	result := MyStorage.AddData(db, cookie, data)

	var retrievedData []byte
	err = db.View(func(tx *bolt.Tx) error {
		bucket := tx.Bucket([]byte("UsersData"))
		if bucket == nil {
			err := errors.New("Bucket not exist")
			return err
		}

		err := bucket.ForEach(func(key, value []byte) error {
			var cookieID MyLogicServer.IDCookie
			if err := json.Unmarshal(key, &cookieID); err != nil {
				return err
			}

			if cookieID.Cookie == cookie && cookieID.ID == result {
				retrievedData = value
			}
			return nil
		})
		if err != nil {
			return err
		}

		return nil
	})
	if err != nil {
		t.Fatalf("Failed to retrieve data from the database: %s", err)
	}

	if !bytes.Equal(retrievedData, data) {
		t.Errorf("Expected retrieved data %s, got %s", string(data), string(retrievedData))
	}
}
