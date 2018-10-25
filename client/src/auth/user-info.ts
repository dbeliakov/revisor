export class UserInfo {
    public username: string;
    public firstName: string;
    public lastName: string;

    public constructor(json: any) {
        this.username = json.username;
        this.firstName = json.first_name;
        this.lastName = json.last_name;
    }
}
