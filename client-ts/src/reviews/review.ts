import { UserInfo } from '@/auth/user-info';

export default class Review {
    public id: string;
    public name: string;
    public closed: boolean;
    public accepted: boolean;
    public owner: UserInfo;
    public reviewers: UserInfo[];
    public updated: Date;

    public constructor(json: any) {
        this.id = json.id;
        this.name = json.name;
        this.closed = json.closed;
        this.accepted = json.accepted;
        this.owner = new UserInfo(json.owner);
        this.reviewers = [];
        for (const reviewer of json.reviewers) {
            this.reviewers.push(new UserInfo(reviewer));
        }
        this.updated = new Date(json.updated * 1000);
    }
}
