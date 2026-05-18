-- +goose Up
-- +goose StatementBegin
CREATE TABLE test_results
(
    id                 BIGSERIAL PRIMARY KEY,
    account_id         BIGINT    NOT NULL,
    module_id          BIGINT    NOT NULL,
    slide_id           BIGINT    NOT NULL,
    selected_option_id BIGINT    NOT NULL,
    is_right           BOOLEAN   NOT NULL,
    attempts           INT       NOT NULL DEFAULT 1,
    first_attempt_at   TIMESTAMP NOT NULL DEFAULT now(),
    last_attempt_at    TIMESTAMP NOT NULL DEFAULT now(),

    CONSTRAINT fk_test_results_account
        FOREIGN KEY (account_id) REFERENCES accounts (id) ON DELETE CASCADE,
    CONSTRAINT fk_test_results_module
        FOREIGN KEY (module_id) REFERENCES modules (id) ON DELETE CASCADE,
    CONSTRAINT fk_test_results_slide
        FOREIGN KEY (slide_id) REFERENCES slides (id) ON DELETE CASCADE,
    CONSTRAINT uq_test_results_account_module_slide
        UNIQUE (account_id, module_id, slide_id)
);

CREATE INDEX idx_test_results_account_id ON test_results (account_id);
CREATE INDEX idx_test_results_slide_id ON test_results (slide_id);
CREATE INDEX idx_test_results_last_attempt_at ON test_results (last_attempt_at DESC);
-- +goose StatementEnd

-- +goose Down
DROP TABLE IF EXISTS test_results;
