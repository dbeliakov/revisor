export class UserInfo {
    public username: string;
    public firstName: string;
    public lastName: string;
    public tgUsername: string|null;

    public constructor(json: any) {
        this.username = json.username;
        this.firstName = json.first_name;
        this.lastName = json.last_name;
        if (json.tg_username) {
            this.tgUsername = json.tg_username;
        } else {
            this.tgUsername = null;
        }
    }
}
