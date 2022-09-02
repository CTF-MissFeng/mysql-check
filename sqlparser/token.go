/*
Copyright 2017 Google Inc.

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

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/xelabs/go-mysqlstack/sqlparser/depends/bytes2"
	"github.com/xelabs/go-mysqlstack/sqlparser/depends/sqltypes"
)

const eofChar = 0x100

// Tokenizer is the struct used to generate SQL
// tokens for the parser.
type Tokenizer struct {
	InStream      *strings.Reader
	AllowComments bool
	ForceEOF      bool
	lastChar      uint16
	Position      int
	lastToken     []byte
	LastError     string
	posVarIndex   int
	ParseTree     Statement
	partialDDL    *DDL
	nesting       int
}

// NewStringTokenizer creates a new Tokenizer for the
// sql string.
func NewStringTokenizer(sql string) *Tokenizer {
	return &Tokenizer{InStream: strings.NewReader(sql)}
}

// keywords is a map of mysql keywords that fall into two categories:
// 1) keywords considered reserved by MySQL
// 2) keywords for us to handle specially in sql.y
//
// Those marked as UNUSED are likely reserved keywords. We add them here so that
// when rewriting queries we can properly backtick quote them so they don't cause issues
//
// NOTE: If you add new keywords, add them also to the reserved_keywords or
// non_reserved_keywords grammar in sql.y -- this will allow the keyword to be used
// in identifiers. See the docs for each grammar to determine which one to put it into.
var keywords = map[string]int{
	"accessible":          UNUSED,
	"add":                 ADD,
	"against":             AGAINST,
	"algorithm":           ALGORITHM,
	"all":                 ALL,
	"alter":               ALTER,
	"analyze":             ANALYZE,
	"and":                 AND,
	"as":                  AS,
	"asc":                 ASC,
	"asensitive":          UNUSED,
	"attach":              ATTACH,
	"attachlist":          ATTACHLIST,
	"auto_increment":      AUTO_INCREMENT,
	"avg_row_length":      AVG_ROW_LENGTH,
	"before":              UNUSED,
	"begin":               BEGIN,
	"between":             BETWEEN,
	"bigint":              BIGINT,
	"binary":              BINARY,
	"binlog":              BINLOG,
	"bit":                 BIT,
	"blob":                BLOB,
	"bool":                BOOL,
	"boolean":             BOOLEAN,
	"both":                UNUSED,
	"btree":               BTREE,
	"by":                  BY,
	"call":                UNUSED,
	"cascade":             UNUSED,
	"case":                CASE,
	"cast":                CAST,
	"change":              UNUSED,
	"char":                CHAR,
	"character":           CHARACTER,
	"charset":             CHARSET,
	"check":               UNUSED,
	"checksum":            CHECKSUM,
	"cleanup":             CLEANUP,
	"collate":             COLLATE,
	"column":              COLUMN,
	"column_format":       COLUMN_FORMAT,
	"columns":             COLUMNS,
	"comment":             COMMENT_KEYWORD,
	"commit":              COMMIT,
	"committed":           COMMITTED,
	"compact":             COMPACT,
	"compression":         COMPRESSION,
	"compressed":          COMPRESSED,
	"connection":          CONNECTION,
	"condition":           UNUSED,
	"constraint":          UNUSED,
	"continue":            UNUSED,
	"convert":             CONVERT,
	"create":              CREATE,
	"cross":               CROSS,
	"current_date":        CURRENT_DATE,
	"current_time":        CURRENT_TIME,
	"current_timestamp":   CURRENT_TIMESTAMP,
	"current_user":        UNUSED,
	"cursor":              UNUSED,
	"data":                DATA,
	"database":            DATABASE,
	"databases":           DATABASES,
	"day_hour":            UNUSED,
	"day_microsecond":     UNUSED,
	"day_minute":          UNUSED,
	"day_second":          UNUSED,
	"date":                DATE,
	"datetime":            DATETIME,
	"dec":                 UNUSED,
	"decimal":             DECIMAL,
	"declare":             UNUSED,
	"default":             DEFAULT,
	"delayed":             UNUSED,
	"delete":              DELETE,
	"desc":                DESC,
	"describe":            DESCRIBE,
	"delay_key_write":     DELAY_KEY_WRITE,
	"detach":              DETACH,
	"deterministic":       UNUSED,
	"directory":           DIRECTORY,
	"distinct":            DISTINCT,
	"distinctrow":         UNUSED,
	"div":                 DIV,
	"distributed":         DISTRIBUTED,
	"disk":                DISK,
	"double":              DOUBLE,
	"drop":                DROP,
	"duplicate":           DUPLICATE,
	"dynamic":             DYNAMIC,
	"each":                UNUSED,
	"else":                ELSE,
	"elseif":              UNUSED,
	"enclosed":            UNUSED,
	"encryption":          ENCRYPTION,
	"end":                 END,
	"engine":              ENGINE,
	"engines":             ENGINES,
	"enum":                ENUM,
	"escape":              ESCAPE,
	"escaped":             UNUSED,
	"events":              EVENTS,
	"exists":              EXISTS,
	"exit":                UNUSED,
	"explain":             EXPLAIN,
	"expansion":           EXPANSION,
	"false":               FALSE,
	"fetch":               UNUSED,
	"fields":              FIELDS,
	"fixed":               FIXED,
	"float":               FLOAT_TYPE,
	"float4":              UNUSED,
	"float8":              UNUSED,
	"for":                 FOR,
	"force":               FORCE,
	"foreign":             UNUSED,
	"from":                FROM,
	"full":                FULL,
	"fulltext":            FULLTEXT,
	"generated":           UNUSED,
	"geometry":            GEOMETRY,
	"geometrycollection":  GEOMETRYCOLLECTION,
	"get":                 UNUSED,
	"global":              GLOBAL,
	"grant":               UNUSED,
	"group":               GROUP,
	"group_concat":        GROUP_CONCAT,
	"gtid":                GTID,
	"hash":                HASH,
	"having":              HAVING,
	"high_priority":       UNUSED,
	"hour_microsecond":    UNUSED,
	"hour_minute":         UNUSED,
	"hour_second":         UNUSED,
	"if":                  IF,
	"ignore":              IGNORE,
	"in":                  IN,
	"index":               INDEX,
	"infile":              UNUSED,
	"inout":               UNUSED,
	"inner":               INNER,
	"insensitive":         UNUSED,
	"insert":              INSERT,
	"insert_method":       INSERT_METHOD,
	"int":                 INT,
	"int1":                UNUSED,
	"int2":                UNUSED,
	"int3":                UNUSED,
	"int4":                UNUSED,
	"int8":                UNUSED,
	"integer":             INTEGER,
	"interval":            INTERVAL,
	"into":                INTO,
	"io_after_gtids":      UNUSED,
	"is":                  IS,
	"isolation":           ISOLATION,
	"iterate":             UNUSED,
	"join":                JOIN,
	"json":                JSON,
	"key":                 KEY,
	"keys":                UNUSED,
	"key_block_size":      KEY_BLOCK_SIZE,
	"kill":                KILL,
	"language":            LANGUAGE,
	"last_insert_id":      LAST_INSERT_ID,
	"leading":             UNUSED,
	"leave":               UNUSED,
	"left":                LEFT,
	"level":               LEVEL,
	"like":                LIKE,
	"limit":               LIMIT,
	"linear":              UNUSED,
	"lines":               UNUSED,
	"linestring":          LINESTRING,
	"list":                LIST,
	"load":                UNUSED,
	"local":               LOCAL,
	"localtime":           LOCALTIME,
	"localtimestamp":      LOCALTIMESTAMP,
	"lock":                LOCK,
	"long":                UNUSED,
	"longblob":            LONGBLOB,
	"longtext":            LONGTEXT,
	"loop":                UNUSED,
	"low_priority":        UNUSED,
	"max_rows":            MAX_ROWS,
	"master_bind":         UNUSED,
	"match":               MATCH,
	"maxvalue":            UNUSED,
	"mediumblob":          MEDIUMBLOB,
	"mediumint":           MEDIUMINT,
	"mediumtext":          MEDIUMTEXT,
	"memory":              MEMORY,
	"multilinestring":     MULTILINESTRING,
	"multipoint":          MULTIPOINT,
	"multipolygon":        MULTIPOLYGON,
	"names":               NAMES,
	"middleint":           UNUSED,
	"minute_microsecond":  UNUSED,
	"minute_second":       UNUSED,
	"min_rows":            MIN_ROWS,
	"mod":                 MOD,
	"mode":                MODE,
	"modify":              MODIFY,
	"modifies":            UNUSED,
	"natural":             NATURAL,
	"nchar":               NCHAR,
	"next":                NEXT,
	"ngram":               NGRAM,
	"not":                 NOT,
	"no_write_to_binlog":  UNUSED,
	"null":                NULL,
	"numeric":             NUMERIC,
	"offset":              OFFSET,
	"on":                  ON,
	"only":                ONLY,
	"optimize":            OPTIMIZE,
	"optimizer_costs":     UNUSED,
	"option":              UNUSED,
	"optionally":          UNUSED,
	"or":                  OR,
	"order":               ORDER,
	"out":                 UNUSED,
	"outer":               OUTER,
	"outfile":             UNUSED,
	"parser":              PARSER,
	"partition":           PARTITION,
	"partitions":          PARTITIONS,
	"pack_keys":           PACK_KEYS,
	"password":            PASSWORD,
	"point":               POINT,
	"polygon":             POLYGON,
	"precision":           UNUSED,
	"primary":             PRIMARY,
	"procedure":           UNUSED,
	"processlist":         PROCESSLIST,
	"query":               QUERY,
	"queryz":              QUERYZ,
	"radon":               RADON,
	"range":               UNUSED,
	"read":                READ,
	"reads":               UNUSED,
	"read_write":          UNUSED,
	"real":                REAL,
	"references":          UNUSED,
	"regexp":              REGEXP,
	"release":             UNUSED,
	"rename":              RENAME,
	"repair":              REPAIR,
	"repeat":              UNUSED,
	"repeatable":          REPEATABLE,
	"replace":             REPLACE,
	"require":             UNUSED,
	"reshard":             RESHARD,
	"resignal":            UNUSED,
	"restrict":            UNUSED,
	"return":              UNUSED,
	"revoke":              UNUSED,
	"right":               RIGHT,
	"rlike":               REGEXP,
	"redundant":           REDUNDANT,
	"rollback":            ROLLBACK,
	"schema":              UNUSED,
	"schemas":             UNUSED,
	"second_microsecond":  UNUSED,
	"select":              SELECT,
	"sensitive":           UNUSED,
	"separator":           SEPARATOR,
	"serializable":        SERIALIZABLE,
	"session":             SESSION,
	"set":                 SET,
	"share":               SHARE,
	"show":                SHOW,
	"signal":              UNUSED,
	"signed":              SIGNED,
	"single":              SINGLE,
	"smallint":            SMALLINT,
	"spatial":             SPATIAL,
	"specific":            UNUSED,
	"sql":                 UNUSED,
	"sqlexception":        UNUSED,
	"sqlstate":            UNUSED,
	"sqlwarning":          UNUSED,
	"sql_big_result":      UNUSED,
	"sql_cache":           SQL_CACHE,
	"sql_calc_found_rows": UNUSED,
	"sql_no_cache":        SQL_NO_CACHE,
	"sql_small_result":    UNUSED,
	"ssl":                 UNUSED,
	"row_format":          ROW_FORMAT,
	"status":              STATUS,
	"start":               START,
	"starting":            UNUSED,
	"stored":              UNUSED,
	"storage":             STORAGE,
	"straight_join":       STRAIGHT_JOIN,
	"stats_auto_recalc":   STATS_AUTO_RECALC,
	"stats_persistent":    STATS_PERSISTENT,
	"stats_sample_pages":  STATS_SAMPLE_PAGES,
	"table":               TABLE,
	"tables":              TABLES,
	"tablespace":          TABLESPACE,
	"terminated":          UNUSED,
	"text":                TEXT,
	"then":                THEN,
	"time":                TIME,
	"timestamp":           TIMESTAMP,
	"tinyblob":            TINYBLOB,
	"tinyint":             TINYINT,
	"tinytext":            TINYTEXT,
	"to":                  TO,
	"tokudb_default":      TOKUDB_DEFAULT,
	"tokudb_fast":         TOKUDB_FAST,
	"tokudb_small":        TOKUDB_SMALL,
	"tokudb_zlib":         TOKUDB_ZLIB,
	"tokudb_quicklz":      TOKUDB_QUICKLZ,
	"tokudb_lzma":         TOKUDB_LZMA,
	"tokudb_snappy":       TOKUDB_SNAPPY,
	"tokudb_uncompressed": TOKUDB_UNCOMPRESSED,
	"trailing":            UNUSED,
	"trigger":             UNUSED,
	"true":                TRUE,
	"truncate":            TRUNCATE,
	"transaction":         TRANSACTION,
	"txnz":                TXNZ,
	"uncommitted":         UNCOMMITTED,
	"undo":                UNUSED,
	"union":               UNION,
	"unique":              UNIQUE,
	"unlock":              UNUSED,
	"unsigned":            UNSIGNED,
	"update":              UPDATE,
	"usage":               UNUSED,
	"use":                 USE,
	"using":               USING,
	"utc_date":            UTC_DATE,
	"utc_time":            UTC_TIME,
	"utc_timestamp":       UTC_TIMESTAMP,
	"values":              VALUES,
	"varbinary":           VARBINARY,
	"varchar":             VARCHAR,
	"varcharacter":        UNUSED,
	"variables":           VARIABLES,
	"varying":             UNUSED,
	"versions":            VERSIONS,
	"virtual":             UNUSED,
	"view":                VIEW,
	"warnings":            WARNINGS,
	"when":                WHEN,
	"where":               WHERE,
	"while":               UNUSED,
	"with":                WITH,
	"write":               WRITE,
	"xa":                  XA,
	"xor":                 UNUSED,
	"year":                YEAR,
	"year_month":          UNUSED,
	"zerofill":            ZEROFILL,
}

// keywordStrings contains the reverse mapping of token to keyword strings
var keywordStrings = map[int]string{}

func init() {
	for str, id := range keywords {
		if id == UNUSED {
			continue
		}
		keywordStrings[id] = str
	}
}

// Lex returns the next token form the Tokenizer.
// This function is used by go yacc.
func (tkn *Tokenizer) Lex(lval *yySymType) int {
	typ, val := tkn.Scan()
	for typ == COMMENT {
		if tkn.AllowComments {
			break
		}
		typ, val = tkn.Scan()
	}
	lval.bytes = val
	tkn.lastToken = val
	return typ
}

// Error is called by go yacc if there's a parsing error.
func (tkn *Tokenizer) Error(err string) {
	buf := &bytes2.Buffer{}
	if tkn.lastToken != nil {
		fmt.Fprintf(buf, "%s at position %v near '%s'", err, tkn.Position, tkn.lastToken)
	} else {
		fmt.Fprintf(buf, "%s at position %v", err, tkn.Position)
	}
	tkn.LastError = buf.String()
}

// Scan scans the tokenizer for the next token and returns
// the token type and an optional value.
func (tkn *Tokenizer) Scan() (int, []byte) {
	if tkn.ForceEOF {
		return 0, nil
	}

	if tkn.lastChar == 0 {
		tkn.next()
	}
	tkn.skipBlank()
	switch ch := tkn.lastChar; {
	case isLetter(ch):
		tkn.next()
		if ch == 'X' || ch == 'x' {
			if tkn.lastChar == '\'' {
				tkn.next()
				return tkn.scanHex()
			}
		}
		isDbSystemVariable := false
		if ch == '@' && tkn.lastChar == '@' {
			isDbSystemVariable = true
		}
		return tkn.scanIdentifier(byte(ch), isDbSystemVariable)
	case isDigit(ch):
		return tkn.scanNumber(false)
	case ch == ':':
		return tkn.scanBindVar()
	default:
		tkn.next()
		switch ch {
		case eofChar:
			return 0, nil
		case '=', ',', ';', '(', ')', '+', '*', '%', '^', '~':
			return int(ch), nil
		case '&':
			if tkn.lastChar == '&' {
				tkn.next()
				return AND, nil
			}
			return int(ch), nil
		case '|':
			if tkn.lastChar == '|' {
				tkn.next()
				return OR, nil
			}
			return int(ch), nil
		case '?':
			tkn.posVarIndex++
			buf := new(bytes2.Buffer)
			fmt.Fprintf(buf, ":v%d", tkn.posVarIndex)
			return VALUE_ARG, buf.Bytes()
		case '.':
			if isDigit(tkn.lastChar) {
				return tkn.scanNumber(true)
			}
			return int(ch), nil
		case '/':
			switch tkn.lastChar {
			case '/':
				tkn.next()
				return tkn.scanCommentType1("//")
			case '*':
				tkn.next()
				return tkn.scanCommentType2()
			default:
				return int(ch), nil
			}
		case '#':
			tkn.next()
			return tkn.scanCommentType1("#")
		case '-':
			switch tkn.lastChar {
			case '-':
				tkn.next()
				return tkn.scanCommentType1("--")
			case '>':
				tkn.next()
				if tkn.lastChar == '>' {
					tkn.next()
					return JSON_UNQUOTE_EXTRACT_OP, nil
				}
				return JSON_EXTRACT_OP, nil
			}
			return int(ch), nil
		case '<':
			switch tkn.lastChar {
			case '>':
				tkn.next()
				return NE, nil
			case '<':
				tkn.next()
				return SHIFT_LEFT, nil
			case '=':
				tkn.next()
				switch tkn.lastChar {
				case '>':
					tkn.next()
					return NULL_SAFE_EQUAL, nil
				default:
					return LE, nil
				}
			default:
				return int(ch), nil
			}
		case '>':
			switch tkn.lastChar {
			case '=':
				tkn.next()
				return GE, nil
			case '>':
				tkn.next()
				return SHIFT_RIGHT, nil
			default:
				return int(ch), nil
			}
		case '!':
			if tkn.lastChar == '=' {
				tkn.next()
				return NE, nil
			}
			return int(ch), nil
		case '\'', '"':
			return tkn.scanString(ch, STRING)
		case '`':
			return tkn.scanLiteralIdentifier()
		default:
			return LEX_ERROR, []byte{byte(ch)}
		}
	}
}

func (tkn *Tokenizer) skipBlank() {
	ch := tkn.lastChar
	for ch == ' ' || ch == '\n' || ch == '\r' || ch == '\t' {
		tkn.next()
		ch = tkn.lastChar
	}
}

func (tkn *Tokenizer) scanIdentifier(firstByte byte, isDbSystemVariable bool) (int, []byte) {
	buffer := &bytes2.Buffer{}
	buffer.WriteByte(firstByte)
	for isLetter(tkn.lastChar) || isDigit(tkn.lastChar) || (isDbSystemVariable && isCarat(tkn.lastChar)) {
		buffer.WriteByte(byte(tkn.lastChar))
		tkn.next()
	}
	lowered := bytes.ToLower(buffer.Bytes())
	loweredStr := string(lowered)
	if keywordID, found := keywords[loweredStr]; found {
		return keywordID, lowered
	}
	// dual must always be case-insensitive
	if loweredStr == "dual" {
		return ID, lowered
	}
	return ID, buffer.Bytes()
}

func isCarat(ch uint16) bool {
	return ch == '.' || ch == '\'' || ch == '"' || ch == '`'
}

func (tkn *Tokenizer) scanHex() (int, []byte) {
	buffer := &bytes2.Buffer{}
	tkn.scanMantissa(16, buffer)
	if tkn.lastChar != '\'' {
		return LEX_ERROR, buffer.Bytes()
	}
	tkn.next()
	if buffer.Len()%2 != 0 {
		return LEX_ERROR, buffer.Bytes()
	}
	return HEX, buffer.Bytes()
}

func (tkn *Tokenizer) scanLiteralIdentifier() (int, []byte) {
	buffer := &bytes2.Buffer{}
	backTickSeen := false
	for {
		if backTickSeen {
			if tkn.lastChar != '`' {
				break
			}
			backTickSeen = false
			buffer.WriteByte('`')
			tkn.next()
			continue
		}
		// The previous char was not a backtick.
		switch tkn.lastChar {
		case '`':
			backTickSeen = true
		case eofChar:
			// Premature EOF.
			return LEX_ERROR, buffer.Bytes()
		default:
			buffer.WriteByte(byte(tkn.lastChar))
		}
		tkn.next()
	}
	if buffer.Len() == 0 {
		return LEX_ERROR, buffer.Bytes()
	}
	return ID, buffer.Bytes()
}

func (tkn *Tokenizer) scanBindVar() (int, []byte) {
	buffer := &bytes2.Buffer{}
	buffer.WriteByte(byte(tkn.lastChar))
	token := VALUE_ARG
	tkn.next()
	if tkn.lastChar == ':' {
		token = LIST_ARG
		buffer.WriteByte(byte(tkn.lastChar))
		tkn.next()
	}
	if !isLetter(tkn.lastChar) {
		return LEX_ERROR, buffer.Bytes()
	}
	for isLetter(tkn.lastChar) || isDigit(tkn.lastChar) || tkn.lastChar == '.' {
		buffer.WriteByte(byte(tkn.lastChar))
		tkn.next()
	}
	return token, buffer.Bytes()
}

func (tkn *Tokenizer) scanMantissa(base int, buffer *bytes2.Buffer) {
	for digitVal(tkn.lastChar) < base {
		tkn.consumeNext(buffer)
	}
}

func (tkn *Tokenizer) scanNumber(seenDecimalPoint bool) (int, []byte) {
	token := INTEGRAL
	buffer := &bytes2.Buffer{}
	if seenDecimalPoint {
		token = FLOAT
		buffer.WriteByte('.')
		tkn.scanMantissa(10, buffer)
		goto exponent
	}

	// 0x construct.
	if tkn.lastChar == '0' {
		tkn.consumeNext(buffer)
		if tkn.lastChar == 'x' || tkn.lastChar == 'X' {
			token = HEXNUM
			tkn.consumeNext(buffer)
			tkn.scanMantissa(16, buffer)
			goto exit
		}
	}

	tkn.scanMantissa(10, buffer)

	if tkn.lastChar == '.' {
		token = FLOAT
		tkn.consumeNext(buffer)
		tkn.scanMantissa(10, buffer)
	}

exponent:
	if tkn.lastChar == 'e' || tkn.lastChar == 'E' {
		token = FLOAT
		tkn.consumeNext(buffer)
		if tkn.lastChar == '+' || tkn.lastChar == '-' {
			tkn.consumeNext(buffer)
		}
		tkn.scanMantissa(10, buffer)
	}

exit:
	// A letter cannot immediately follow a number.
	if isLetter(tkn.lastChar) {
		return LEX_ERROR, buffer.Bytes()
	}

	return token, buffer.Bytes()
}

func (tkn *Tokenizer) scanString(delim uint16, typ int) (int, []byte) {
	buffer := &bytes2.Buffer{}
	for {
		ch := tkn.lastChar
		tkn.next()
		if ch == delim {
			if tkn.lastChar == delim {
				tkn.next()
			} else {
				break
			}
		} else if ch == '\\' {
			if tkn.lastChar == eofChar {
				return LEX_ERROR, buffer.Bytes()
			}
			if decodedChar := sqltypes.SQLDecodeMap[byte(tkn.lastChar)]; decodedChar == sqltypes.DontEscape {
				ch = tkn.lastChar
			} else {
				ch = uint16(decodedChar)
			}
			tkn.next()
		}
		if ch == eofChar {
			return LEX_ERROR, buffer.Bytes()
		}
		buffer.WriteByte(byte(ch))
	}
	return typ, buffer.Bytes()
}

func (tkn *Tokenizer) scanCommentType1(prefix string) (int, []byte) {
	buffer := &bytes2.Buffer{}
	buffer.WriteString(prefix)
	for tkn.lastChar != eofChar {
		if tkn.lastChar == '\n' {
			tkn.consumeNext(buffer)
			break
		}
		tkn.consumeNext(buffer)
	}
	return COMMENT, buffer.Bytes()
}

func (tkn *Tokenizer) scanCommentType2() (int, []byte) {
	buffer := &bytes2.Buffer{}
	buffer.WriteString("/*")
	for {
		if tkn.lastChar == '*' {
			tkn.consumeNext(buffer)
			if tkn.lastChar == '/' {
				tkn.consumeNext(buffer)
				break
			}
			continue
		}
		if tkn.lastChar == eofChar {
			return LEX_ERROR, buffer.Bytes()
		}
		tkn.consumeNext(buffer)
	}
	return COMMENT, buffer.Bytes()
}

func (tkn *Tokenizer) consumeNext(buffer *bytes2.Buffer) {
	if tkn.lastChar == eofChar {
		// This should never happen.
		panic("unexpected EOF")
	}
	buffer.WriteByte(byte(tkn.lastChar))
	tkn.next()
}

func (tkn *Tokenizer) next() {
	if ch, err := tkn.InStream.ReadByte(); err != nil {
		// Only EOF is possible.
		tkn.lastChar = eofChar
	} else {
		tkn.lastChar = uint16(ch)
	}
	tkn.Position++
}

func isLetter(ch uint16) bool {
	return 'a' <= ch && ch <= 'z' || 'A' <= ch && ch <= 'Z' || ch == '_' || ch == '@'
}

func digitVal(ch uint16) int {
	switch {
	case '0' <= ch && ch <= '9':
		return int(ch) - '0'
	case 'a' <= ch && ch <= 'f':
		return int(ch) - 'a' + 10
	case 'A' <= ch && ch <= 'F':
		return int(ch) - 'A' + 10
	}
	return 16 // larger than any legal digit val
}

func isDigit(ch uint16) bool {
	return '0' <= ch && ch <= '9'
}
