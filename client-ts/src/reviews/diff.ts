import Comment from '@/reviews/comment';

class Line {
    public content: string;
    public revision: number;
    public id: string;

    public constructor(json: any) {
        this.content = json.content;
        this.revision = json.revision;
        this.id = json.id;
    }
}

class DiffLine {
    public type: string;
    public old: Line;
    public new: Line;

    public constructor(json: any) {
        this.type = json.type;
        this.old = new Line(json.old);
        this.new = new Line(json.new);
    }
}

class DiffRange {
    public from: number;
    public to: number;

    public constructor(json: any) {
        this.from = json.from;
        this.to = json.to;
    }
}

class DiffGroup {
    public oldRange: DiffRange;
    public newRange: DiffRange;
    public lines: DiffLine[];

    public constructor(json: any) {
        this.oldRange = new DiffRange(json.old_range);
        this.newRange = new DiffRange(json.new_range);
        this.lines = [];
        for (const line of json.lines) {
            this.lines.push(new DiffLine(line));
        }
    }
}

export class Diff {
    public filename: string;
    public groups: DiffGroup[];

    public constructor(json: any) {
        this.filename = json.filename;
        this.groups = []
        for (const group of json.groups) {
            this.groups.push(new DiffGroup(group));
        }
    }
}