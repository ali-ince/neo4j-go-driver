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

package drivertest

import (
	"fmt"
	"github.com/neo4j-drivers/neo4j-go-connector"
	"github.com/onsi/gomega"
	"github.com/onsi/gomega/types"
)

func BeDatabaseError() types.GomegaMatcher {
	return &databaseErrorMatcher{}
}

func BeTransientError() types.GomegaMatcher {
	return &databaseErrorMatcher{
		classificationMatcher: gomega.BeEquivalentTo("TransientError"),
	}
}

func BeClientError() types.GomegaMatcher {
	return &databaseErrorMatcher{
		classificationMatcher: gomega.BeEquivalentTo("ClientError"),
	}
}

func BeArithmeticError() types.GomegaMatcher {
	return &databaseErrorMatcher{
		codeMatcher: gomega.ContainSubstring("ArithmeticError"),
	}
}

func BeSyntaxError() types.GomegaMatcher {
	return &databaseErrorMatcher{
		codeMatcher: gomega.ContainSubstring("SyntaxError"),
	}
}

func BeServiceUnavailableError() types.GomegaMatcher {
	return &serviceUnavailableErrorMatcher{}
}

func BePoolFullError() types.GomegaMatcher {
	return &poolErrorMatcher{
		expected: 1,
		expectedText: "Pool full",
	}
}

func BeAuthenticationError() types.GomegaMatcher {
	return &connectorErrorMatcher{
		stateMatcher: gomega.BeEquivalentTo(4),
		codeMatcher: gomega.BeEquivalentTo(7),
	}
}

func ContainMessage(part string) types.GomegaMatcher {
	return &databaseErrorMatcher{
		messageMatcher: gomega.ContainSubstring(part),
	}
}

type databaseErrorMatcher struct {
	classificationMatcher types.GomegaMatcher
	codeMatcher           types.GomegaMatcher
	messageMatcher        types.GomegaMatcher
}

type serviceUnavailableErrorMatcher struct {
}

type poolErrorMatcher struct {
	expected uint32
	expectedText string
}

type connectorErrorMatcher struct {
	stateMatcher types.GomegaMatcher
	codeMatcher  types.GomegaMatcher
}

func (matcher *databaseErrorMatcher) Match(actual interface{}) (success bool, err error) {
	databaseError, ok := actual.(*seabolt.DatabaseError)
	if !ok {
		return false, nil
	}

	if matcher.classificationMatcher != nil {
		return matcher.classificationMatcher.Match(databaseError.Classification())
	}

	if matcher.codeMatcher != nil {
		return matcher.codeMatcher.Match(databaseError.Code())
	}

	if matcher.messageMatcher != nil {
		return matcher.messageMatcher.Match(databaseError.Message())
	}

	return true, nil
}

func (matcher *databaseErrorMatcher) FailureMessage(actual interface{}) (message string) {
	databaseError, ok := actual.(*seabolt.DatabaseError)
	if !ok {
		return fmt.Sprintf("Expected\n\t%#v\nto be a DatabaseError", actual)
	}

	if matcher.classificationMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nto have its classification to match %s", actual, matcher.classificationMatcher.FailureMessage(databaseError.Classification()))
	}

	if matcher.codeMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nto have its code to match %s", actual, matcher.codeMatcher.FailureMessage(databaseError.Code()))
	}

	if matcher.messageMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nto have its message to match %s", actual, matcher.classificationMatcher.FailureMessage(databaseError.Message()))
	}

	return fmt.Sprintf("Unexpected condition in matcher")
}

func (matcher *databaseErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	databaseError, ok := actual.(*seabolt.DatabaseError)
	if !ok {
		return fmt.Sprintf("Expected\n\t%#v\nnot to be a DatabaseError", actual)
	}

	if matcher.classificationMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nnot to have its classification to match %s", actual, matcher.classificationMatcher.FailureMessage(databaseError.Classification()))
	}

	if matcher.codeMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nnot to have its code to match %s", actual, matcher.codeMatcher.FailureMessage(databaseError.Code()))
	}

	if matcher.messageMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nnot to have its message to match %s", actual, matcher.classificationMatcher.FailureMessage(databaseError.Message()))
	}

	return fmt.Sprintf("Unexpected condition in matcher")
}

func (matcher *serviceUnavailableErrorMatcher) Match(actual interface{}) (success bool, err error) {
	err, ok := actual.(error)
	if !ok {
		return false, nil
	}

	return seabolt.IsServiceUnavailable(err), nil
}

func (matcher *serviceUnavailableErrorMatcher) FailureMessage(actual interface{}) (message string) {
	_, ok := actual.(error)
	if !ok {
		return fmt.Sprintf("Expected\n\t%#v\nto be an error", actual)
	}

	return fmt.Sprintf("Expected\n\t%#v\nto be a ServiceUnavailableError", actual)
}

func (matcher *serviceUnavailableErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	_, ok := actual.(error)
	if !ok {
		return fmt.Sprintf("Expected\n\t%#v\nnot to be an error", actual)
	}

	return fmt.Sprintf("Expected\n\t%#v\nnot to be a ServiceUnavailableError", actual)
}

func (matcher *poolErrorMatcher) Match(actual interface{}) (success bool, err error) {
	poolError, ok := actual.(*seabolt.PoolError)
	if !ok {
		return false, nil
	}

	return poolError.Code() == matcher.expected, nil
}

func (matcher *poolErrorMatcher) FailureMessage(actual interface{}) (message string) {
	_, ok := actual.(*seabolt.PoolError)
	if !ok {
		return fmt.Sprintf("Expected\n\t%#v\nto be a pool error", actual)
	}

	return fmt.Sprintf("Expected\n\t%#v\nto be a pool error with code: %d (%s)", actual, matcher.expected, matcher.expectedText)
}

func (matcher *poolErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	_, ok := actual.(*seabolt.PoolError)
	if !ok {
		return fmt.Sprintf("Expected\n\t%#v\nnot to be a pool error", actual)
	}

	return fmt.Sprintf("Expected\n\t%#v\nnot to be a pool error with code: %d (%s)", actual, matcher.expected, matcher.expectedText)
}

func (matcher *connectorErrorMatcher) Match(actual interface{}) (success bool, err error) {
	connectorError, ok := actual.(*seabolt.ConnectorError)
	if !ok {
		return false, nil
	}

	if matcher.stateMatcher != nil {
		return matcher.stateMatcher.Match(connectorError.State())
	}

	if matcher.codeMatcher != nil {
		return matcher.codeMatcher.Match(connectorError.Code())
	}

	return true, nil
}

func (matcher *connectorErrorMatcher) FailureMessage(actual interface{}) (message string) {
	connectorError, ok := actual.(*seabolt.ConnectorError)
	if !ok {
		return fmt.Sprintf("Expected\n\t%#v\nto be a ConnectorError", actual)
	}

	if matcher.stateMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nto have its state to match %s", actual, matcher.stateMatcher.FailureMessage(connectorError.State()))
	}

	if matcher.codeMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nto have its code to match %s", actual, matcher.codeMatcher.FailureMessage(connectorError.Code()))
	}

	return fmt.Sprintf("Unexpected condition in matcher")
}

func (matcher *connectorErrorMatcher) NegatedFailureMessage(actual interface{}) (message string) {
	connectorError, ok := actual.(*seabolt.ConnectorError)
	if !ok {
		return fmt.Sprintf("Expected\n\t%#v\nnot to be a ConnectorError", actual)
	}

	if matcher.stateMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nnot to have its state to match %s", actual, matcher.stateMatcher.FailureMessage(connectorError.State()))
	}

	if matcher.codeMatcher != nil {
		return fmt.Sprintf("Expected\n\t%#v\nnot to have its code to match %s", actual, matcher.codeMatcher.FailureMessage(connectorError.Code()))
	}

	return fmt.Sprintf("Unexpected condition in matcher")
}
