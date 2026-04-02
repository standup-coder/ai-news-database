package db

import (
	"database/sql"
	"fmt"
	"news4coder/internal/article"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

const dbFile = "news4coder.db"

// DB 封装数据库操作
type DB struct {
	conn *sql.DB
}

// New 创建并初始化数据库
func New() (*DB, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return nil, fmt.Errorf("无法获取用户主目录: %w", err)
	}

	dbDir := filepath.Join(homeDir, ".news4coder")
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return nil, fmt.Errorf("无法创建配置目录: %w", err)
	}

	dbPath := filepath.Join(dbDir, dbFile)
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, fmt.Errorf("打开数据库失败: %w", err)
	}

	if err := conn.Ping(); err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	db := &DB{conn: conn}
	if err := db.migrate(); err != nil {
		return nil, fmt.Errorf("数据库迁移失败: %w", err)
	}

	return db, nil
}

// Close 关闭数据库连接
func (d *DB) Close() error {
	return d.conn.Close()
}

// migrate 执行数据库初始化
func (d *DB) migrate() error {
	schema := `
CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL,
    url TEXT NOT NULL UNIQUE,
    source TEXT NOT NULL,
    source_alias TEXT NOT NULL,
    summary TEXT,
    llm_summary TEXT,
    llm_tags TEXT,
    quality_score REAL DEFAULT 0,
    language TEXT,
    raw_content TEXT,
    published_at DATETIME,
    fetched_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    enriched_at DATETIME,
    read_status TEXT DEFAULT 'unread',
    tags TEXT,
    note TEXT
);

CREATE INDEX IF NOT EXISTS idx_articles_status ON articles(read_status);
CREATE INDEX IF NOT EXISTS idx_articles_source ON articles(source_alias);
CREATE INDEX IF NOT EXISTS idx_articles_fetched_at ON articles(fetched_at);

CREATE VIRTUAL TABLE IF NOT EXISTS articles_fts USING fts5(
    title, summary, llm_summary, content='articles', content_rowid='id'
);

CREATE TRIGGER IF NOT EXISTS articles_ai AFTER INSERT ON articles BEGIN
    INSERT INTO articles_fts(rowid, title, summary, llm_summary)
    VALUES (new.id, new.title, new.summary, new.llm_summary);
END;

CREATE TRIGGER IF NOT EXISTS articles_ad AFTER DELETE ON articles BEGIN
    INSERT INTO articles_fts(articles_fts, rowid, title, summary, llm_summary)
    VALUES ('delete', old.id, old.title, old.summary, old.llm_summary);
END;

CREATE TRIGGER IF NOT EXISTS articles_au AFTER UPDATE ON articles BEGIN
    INSERT INTO articles_fts(articles_fts, rowid, title, summary, llm_summary)
    VALUES ('delete', old.id, old.title, old.summary, old.llm_summary);
    INSERT INTO articles_fts(rowid, title, summary, llm_summary)
    VALUES (new.id, new.title, new.summary, new.llm_summary);
END;
`
	if _, err := d.conn.Exec(schema); err != nil {
		return err
	}

	// 兼容旧数据库：添加新字段
	newColumns := []string{
		"llm_summary TEXT",
		"llm_tags TEXT",
		"quality_score REAL DEFAULT 0",
		"language TEXT",
		"raw_content TEXT",
		"enriched_at DATETIME",
	}
	for _, col := range newColumns {
		_, _ = d.conn.Exec(fmt.Sprintf("ALTER TABLE articles ADD COLUMN %s", col))
		// 忽略错误（字段已存在时会报错）
	}

	// 将旧数据中 NULL 的 quality_score 设为 0
	_, _ = d.conn.Exec("UPDATE articles SET quality_score = 0 WHERE quality_score IS NULL")

	return nil
}

// SaveArticle 保存或更新文章（URL 唯一）
func (d *DB) SaveArticle(a *article.Article) error {
	// 如果已存在则更新基础字段
	var existingID int64
	err := d.conn.QueryRow("SELECT id FROM articles WHERE url = ?", a.URL).Scan(&existingID)
	if err == nil {
		_, err = d.conn.Exec(
			`UPDATE articles SET title = ?, summary = ?, source = ?, source_alias = ?, raw_content = ?, published_at = ? WHERE id = ?`,
			a.Title, a.Summary, a.Source, a.SourceAlias, a.RawContent, a.PublishedAt, existingID,
		)
		return err
	}

	_, err = d.conn.Exec(
		`INSERT INTO articles (title, url, source, source_alias, summary, raw_content, published_at, fetched_at, read_status, tags, note)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		a.Title, a.URL, a.Source, a.SourceAlias, a.Summary, a.RawContent, a.PublishedAt, a.FetchedAt, a.ReadStatus, a.Tags, a.Note,
	)
	return err
}

// GetArticles 按条件查询文章
func (d *DB) GetArticles(status article.ReadStatus, sourceAlias string, limit int) ([]article.Article, error) {
	query := "SELECT id, title, url, source, source_alias, summary, llm_summary, llm_tags, quality_score, language, raw_content, published_at, fetched_at, enriched_at, read_status, tags, note FROM articles WHERE 1=1"
	args := []any{}

	if status != "" {
		query += " AND read_status = ?"
		args = append(args, status)
	}
	if sourceAlias != "" {
		query += " AND source_alias = ?"
		args = append(args, sourceAlias)
	}
	query += " ORDER BY fetched_at DESC"
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := d.conn.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticles(rows)
}

// SearchArticles 全文搜索文章
func (d *DB) SearchArticles(keyword string, limit int) ([]article.Article, error) {
	query := `
		SELECT a.id, a.title, a.url, a.source, a.source_alias, a.summary, a.llm_summary, a.llm_tags, a.quality_score, a.language, a.raw_content, a.published_at, a.fetched_at, a.enriched_at, a.read_status, a.tags, a.note
		FROM articles_fts fts
		JOIN articles a ON a.id = fts.rowid
		WHERE articles_fts MATCH ?
		ORDER BY rank
	`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}

	rows, err := d.conn.Query(query, keyword)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanArticles(rows)
}

// UpdateStatus 更新文章阅读状态
func (d *DB) UpdateStatus(id int64, status article.ReadStatus) error {
	_, err := d.conn.Exec("UPDATE articles SET read_status = ? WHERE id = ?", status, id)
	return err
}

// AddNote 为文章添加笔记
func (d *DB) AddNote(id int64, note string) error {
	_, err := d.conn.Exec("UPDATE articles SET note = ? WHERE id = ?", note, id)
	return err
}

// AddTags 为文章添加标签
func (d *DB) AddTags(id int64, tags string) error {
	_, err := d.conn.Exec("UPDATE articles SET tags = ? WHERE id = ?", tags, id)
	return err
}

// GetUnenrichedArticles 获取尚未经过 LLM 增强的文章
func (d *DB) GetUnenrichedArticles(limit int) ([]article.Article, error) {
	query := `SELECT id, title, url, source, source_alias, summary, llm_summary, llm_tags, quality_score, language, raw_content, published_at, fetched_at, enriched_at, read_status, tags, note
		FROM articles WHERE enriched_at IS NULL ORDER BY fetched_at DESC`
	if limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", limit)
	}
	rows, err := d.conn.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanArticles(rows)
}

// UpdateEnrichment 更新 LLM 增强结果
func (d *DB) UpdateEnrichment(id int64, llmSummary, llmTags, language string, score float64) error {
	_, err := d.conn.Exec(
		`UPDATE articles SET llm_summary = ?, llm_tags = ?, quality_score = ?, language = ?, enriched_at = CURRENT_TIMESTAMP WHERE id = ?`,
		llmSummary, llmTags, score, language, id,
	)
	return err
}

// UpdateRawContent 更新原始内容
func (d *DB) UpdateRawContent(id int64, rawContent string) error {
	_, err := d.conn.Exec("UPDATE articles SET raw_content = ? WHERE id = ?", rawContent, id)
	return err
}

// SearchByKeyword 使用 LIKE 进行模糊搜索（FTS 的 fallback）
func (d *DB) SearchByKeyword(keyword string, limit int) ([]article.Article, error) {
	query := `SELECT id, title, url, source, source_alias, summary, llm_summary, llm_tags, quality_score, language, raw_content, published_at, fetched_at, enriched_at, read_status, tags, note
		FROM articles
		WHERE title LIKE ? OR summary LIKE ? OR llm_summary LIKE ? OR llm_tags LIKE ?
		ORDER BY quality_score DESC, fetched_at DESC
		LIMIT ?`
	pattern := "%" + keyword + "%"
	rows, err := d.conn.Query(query, pattern, pattern, pattern, pattern, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanArticles(rows)
}

// DeleteArticlesByStatus 按状态删除文章（清理用）
func (d *DB) DeleteArticlesByStatus(status article.ReadStatus, beforeDays int) error {
	// SQLite 参数绑定不能在字符串字面量中使用，需要拼接 SQL
	query := fmt.Sprintf(
		"DELETE FROM articles WHERE read_status = ? AND fetched_at < datetime('now', '-%d days')",
		beforeDays,
	)
	_, err := d.conn.Exec(query, status)
	return err
}

// Stats 返回按源统计的数据
func (d *DB) Stats() (map[string]map[string]any, error) {
	rows, err := d.conn.Query(`
		SELECT source_alias,
		       COUNT(*) as total,
		       SUM(CASE WHEN read_status = 'read' THEN 1 ELSE 0 END) as read_count,
		       SUM(CASE WHEN read_status = 'starred' THEN 1 ELSE 0 END) as starred_count,
		       SUM(CASE WHEN read_status = 'unread' THEN 1 ELSE 0 END) as unread_count
		FROM articles
		GROUP BY source_alias
	`)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stats := make(map[string]map[string]any)
	for rows.Next() {
		var alias string
		var total, readCount, starredCount, unreadCount int
		if err := rows.Scan(&alias, &total, &readCount, &starredCount, &unreadCount); err != nil {
			continue
		}
		stats[alias] = map[string]any{
			"total":    total,
			"read":     readCount,
			"starred":  starredCount,
			"unread":   unreadCount,
			"read_rate": float64(readCount+starredCount) / float64(total),
		}
	}
	return stats, nil
}

func scanArticles(rows *sql.Rows) ([]article.Article, error) {
	var articles []article.Article
	for rows.Next() {
		var a article.Article
		var publishedAt, enrichedAt sql.NullTime
		var llmSummary, llmTags, language, rawContent sql.NullString
		var qualityScore sql.NullFloat64
		err := rows.Scan(
			&a.ID, &a.Title, &a.URL, &a.Source, &a.SourceAlias,
			&a.Summary, &llmSummary, &llmTags, &qualityScore, &language, &rawContent,
			&publishedAt, &a.FetchedAt, &enrichedAt, &a.ReadStatus, &a.Tags, &a.Note,
		)
		if err != nil {
			return nil, err
		}
		if publishedAt.Valid {
			a.PublishedAt = &publishedAt.Time
		}
		if enrichedAt.Valid {
			a.EnrichedAt = &enrichedAt.Time
		}
		if llmSummary.Valid {
			a.LLMSummary = llmSummary.String
		}
		if llmTags.Valid {
			a.LLMTags = llmTags.String
		}
		if language.Valid {
			a.Language = language.String
		}
		if rawContent.Valid {
			a.RawContent = rawContent.String
		}
		if qualityScore.Valid {
			a.QualityScore = qualityScore.Float64
		}
		articles = append(articles, a)
	}
	return articles, rows.Err()
}
