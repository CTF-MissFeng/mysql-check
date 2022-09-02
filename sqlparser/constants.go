/*
Copyright 2019 The Vitess Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package sqlparser

const (
	// Select.Distinct
	DistinctStr      = "distinct "
	StraightJoinHint = "straight_join "

	// Select.Lock
	ForUpdateStr = " for update"
	ShareModeStr = " lock in share mode"

	// Select.Cache
	SQLCacheStr   = "sql_cache "
	SQLNoCacheStr = "sql_no_cache "

	// Union.Type
	UnionStr         = "union"
	UnionAllStr      = "union all"
	UnionDistinctStr = "union distinct"

	// InsertStr represents insert action.
	InsertStr = "insert"
	// ReplaceStr represents replace action.
	ReplaceStr = "replace"

	// Set.Scope or Show.Scope.
	SessionStr = "session"
	GlobalStr  = "global"

	// DDL strings.
	CreateDBStr             = "create database"
	CreateTableStr          = "create table"
	CreatePartitionTableStr = "create partition table"
	CreateIndexStr          = "create index"
	DropDBStr               = "drop database"
	DropTableStr            = "drop table"
	DropIndexStr            = "drop index"
	AlterStr                = "alter"
	AlterEngineStr          = "alter table"
	AlterCharsetStr         = "alter table charset"
	AlterAddColumnStr       = "alter table add column"
	AlterDropColumnStr      = "alter table drop column"
	AlterModifyColumnStr    = "alter table modify column"
	RenameStr               = "rename table"
	TruncateTableStr        = "truncate table"
	SingleTableType         = "singletable"
	GlobalTableType         = "globaltable"
	PartitionTableHash      = "partitiontablehash"
	NormalTableType         = "normaltable"
	PartitionTableList      = "partitiontablelist"

	// Index key type strings.
	IndexStr    = "index "
	FullTextStr = "fulltext index "
	SpatialStr  = "spatial index "
	UniqueStr   = "unique index "

	// The following constants represent SHOW statements.
	ShowDatabasesStr      = "databases"
	ShowCreateDatabaseStr = "create database"
	ShowTableStatusStr    = "table status"
	ShowTablesStr         = "tables"
	ShowColumnsStr        = "columns"
	ShowCreateTableStr    = "create table"
	ShowEnginesStr        = "engines"
	ShowStatusStr         = "status"
	ShowVersionsStr       = "versions"
	ShowProcesslistStr    = "processlist"
	ShowQueryzStr         = "queryz"
	ShowTxnzStr           = "txnz"
	ShowWarningsStr       = "warnings"
	ShowVariablesStr      = "variables"
	ShowBinlogEventsStr   = "binlog events"
	ShowUnsupportedStr    = "unsupported"

	// JoinTableExpr.Join.
	JoinStr             = "join"
	StraightJoinStr     = "straight_join"
	LeftJoinStr         = "left join"
	RightJoinStr        = "right join"
	NaturalJoinStr      = "natural join"
	NaturalLeftJoinStr  = "natural left join"
	NaturalRightJoinStr = "natural right join"

	// Index hints.
	UseStr    = "use "
	IgnoreStr = "ignore "
	ForceStr  = "force "

	// Where.Type
	WhereStr  = "where"
	HavingStr = "having"

	// ComparisonExpr.Operator
	EqualStr             = "="
	LessThanStr          = "<"
	GreaterThanStr       = ">"
	LessEqualStr         = "<="
	GreaterEqualStr      = ">="
	NotEqualStr          = "!="
	NullSafeEqualStr     = "<=>"
	InStr                = "in"
	NotInStr             = "not in"
	LikeStr              = "like"
	NotLikeStr           = "not like"
	RegexpStr            = "regexp"
	NotRegexpStr         = "not regexp"
	JSONExtractOp        = "->"
	JSONUnquoteExtractOp = "->>"

	// RangeCond.Operator
	BetweenStr    = "between"
	NotBetweenStr = "not between"

	// IsExpr.Operator
	IsNullStr     = "is null"
	IsNotNullStr  = "is not null"
	IsTrueStr     = "is true"
	IsNotTrueStr  = "is not true"
	IsFalseStr    = "is false"
	IsNotFalseStr = "is not false"

	// BinaryExpr.Operator
	BitAndStr     = "&"
	BitOrStr      = "|"
	BitXorStr     = "^"
	PlusStr       = "+"
	MinusStr      = "-"
	MultStr       = "*"
	DivStr        = "/"
	IntDivStr     = "div"
	ModStr        = "%"
	ShiftLeftStr  = "<<"
	ShiftRightStr = ">>"

	// UnaryExpr.Operator.
	UPlusStr  = "+"
	UMinusStr = "-"
	TildaStr  = "~"
	BangStr   = "!"
	BinaryStr = "binary "

	// this string is "character set" and this comment is required.
	CharacterSetStr = " character set"

	// MatchExpr.Option.
	BooleanModeStr                           = " in boolean mode"
	NaturalLanguageModeStr                   = " in natural language mode"
	NaturalLanguageModeWithQueryExpansionStr = " in natural language mode with query expansion"
	QueryExpansionStr                        = " with query expansion"

	// Order.Direction.
	AscScr  = "asc"
	DescScr = "desc"

	AttachStr     = "attach"
	DetachStr     = "detach"
	AttachListStr = "attachlist"
	ReshardStr    = "reshard"
	CleanupStr    = "cleanup"

	// Transaction isolation levels.
	ReadUncommitted = "read uncommitted"
	ReadCommitted   = "read committed"
	RepeatableRead  = "repeatable read"
	Serializable    = "serializable"

	// Transaction access mode.
	TxReadOnly  = "read only"
	TxReadWrite = "read write"

	// StartTxnStr represents the txn start transaction.
	StartTxnStr = "start transaction"

	// BeginTxnStr represents the txn begin.
	BeginTxnStr = "begin"

	// RollbackTxnStr represents the txn rollback.
	RollbackTxnStr = "rollback"

	// CommitTxnStr represents the txn commit.
	CommitTxnStr = "commit"
)
