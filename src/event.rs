use serde::{Deserialize, Serialize};

#[derive(Default, Debug, Clone, PartialEq, Serialize, Deserialize)]
pub struct Event {
    pub event: String,
    pub sid: String,
    pub thread: String,
    pub time: Option<String>,
    pub evt: Option<String>,
    pub exe: Option<String>,
    pub t_abs: Option<f64>,
    #[serde(default)]
    pub argv: Vec<String>,
    pub repo: Option<i64>,
    pub worktree: Option<String>,
    pub name: Option<String>,
    pub hierarchy: Option<String>,
    pub code: Option<i64>,
    pub nesting: Option<i64>,
    pub category: Option<String>,
    pub label: Option<String>,
    pub msg: Option<String>,
    pub t_rel: Option<f64>,
    pub key: Option<String>,
    pub value: Option<String>,
    pub child_id: Option<i64>,
    pub child_class: Option<String>,
    pub cd: Option<String>,
    pub use_shell: Option<bool>,
    pub pid: Option<i64>,
}

