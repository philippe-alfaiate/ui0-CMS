package main

import "time"

type AdministrationUser struct {
	Id    int64 `pg:",pk"`
	Name  string
	Email []string
	Phone []string
	Token []string
}

type AdministrationAuth struct {
	Id       int64 `pg:",pk"`
	UserId   int64
	User     *AdministrationUser `pg:"rel:has-one"`
	Hash     string
	Salt     string
	Protocol string
}

type AdministrationSessionToken struct {
	Token      string
	Creation   time.Time
	Expiration time.Time
}

type AdministrationSession struct {
	Id    int64
	User  *AdministrationUser `pg:"rel:has-one"`
	Token []AdministrationSessionToken
}
