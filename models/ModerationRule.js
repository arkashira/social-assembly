class ModerationRule {
  constructor(title, description, severity, version = 1, createdAt = new Date()) {
    this.title = title;
    this.description = description;
    this.severity = severity;
    this.version = version;
    this.createdAt = createdAt;
    this.updatedAt = new Date();
  }

  update(title, description, severity) {
    this.title = title;
    this.description = description;
    this.severity = severity;
    this.version++;
    this.updatedAt = new Date();
  }

  archive() {
    this.archived = true;
    this.updatedAt = new Date();
  }
}