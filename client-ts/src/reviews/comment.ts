import { UserInfo } from '@/auth/user-info';

export default class Comment {
    public id: string;
    public author: UserInfo;
    public created: Date;
    public text: string;
    public lineId: string;

    public constructor(json: any) {
        this.id = json.id;
        this.author = new UserInfo(json.author);
        this.created = new Date(json.created * 1000);
        this.text = json.text;
        this.lineId = json.line_id;
    }
}