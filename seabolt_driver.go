/*
 * Copyright (c) 2002-2018 "Neo4j,"
 * Neo4j Sweden AB [http://neo4j.com]
 *
 * This file is part of Neo4j.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package neo4j

import (
	"net/url"
	"sync/atomic"

	"github.com/neo4j-drivers/neo4j-go-connector"
	"github.com/pkg/errors"
)

type seaboltDriver struct {
	config    *Config
	target    url.URL
	connector seabolt.Connector

	open int32
}

func configToSeaboltConfig(config *Config) *seabolt.Config {
	return &seabolt.Config{
		Encryption:  config.Encrypted,
		Log:         config.Log,
		MaxPoolSize: config.MaxConnectionPoolSize,
		ValueHandlers: []seabolt.ValueHandler{
			&nodeValueHandler{},
			&relationshipValueHandler{},
			&pathValueHandler{},
			&pointValueHandler{},
			&dateValueHandler{},
			&localTimeValueHandler{},
			&offsetTimeValueHandler{},
			&localDateTimeValueHandler{},
			&dateTimeValueHandler{},
			&durationValueHandler{},
		},
	}
}

func newSeaboltDriver(target *url.URL, token AuthToken, config *Config) (*seaboltDriver, error) {
	if config == nil {
		config = defaultConfig()
	}

	connector, err := seabolt.NewConnector(target, token.tokens, configToSeaboltConfig(config))
	if err != nil {
		return nil, err
	}

	driver := seaboltDriver{
		config:    config,
		target:    *target,
		connector: connector,
		open:      1,
	}
	return &driver, nil
}

func assertDriverOpen(driver *seaboltDriver) error {
	if atomic.LoadInt32(&driver.open) == 0 {
		return errors.New("cannot acquire a session on a closed driver")
	}

	return nil
}

func (driver *seaboltDriver) Target() url.URL {
	return driver.target
}

func (driver *seaboltDriver) Session(accessMode AccessMode, bookmarks ...string) (*Session, error) {
	if err := assertDriverOpen(driver); err != nil {
		return nil, err
	}

	return newSession(driver, accessMode, bookmarks), nil
}

func (driver *seaboltDriver) Close() error {
	if atomic.CompareAndSwapInt32(&driver.open, 1, 0) {
		return driver.connector.Close()
	}

	return nil
}

func (driver *seaboltDriver) configuration() *Config {
	return driver.config
}

func (driver *seaboltDriver) acquire(mode AccessMode) (seabolt.Connection, error) {
	if err := assertDriverOpen(driver); err != nil {
		return nil, err
	}

	seaboltMode := seabolt.AccessModeWrite
	if mode == AccessModeRead {
		seaboltMode = seabolt.AccessModeRead
	}

	return driver.connector.Acquire(seaboltMode)
}