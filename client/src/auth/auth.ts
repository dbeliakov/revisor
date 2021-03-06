import {AxiosStatic, AxiosResponse} from 'axios';
import {UserInfo} from '@/auth/user-info';
import {Component, Vue, Watch} from 'vue-property-decorator';
import { resolve } from 'dns';
import { responseToError } from '@/utils/utils';

export default class Auth {
    public axios: AxiosStatic; // TODO make it private
    private watch: any;

    public constructor(axios: AxiosStatic) {
        this.axios = axios;

        this.watch = new Vue({
           data() {
               return {
                   ready: false,
                   authenticated: false,
                   userInfo: undefined,
               };
           },
        });

        const storage = window.localStorage;
        const token = storage.getItem('auth_header_value');
        if (!token) {
            this.watch.ready = true;
            this.watch.authenticated = false;
        } else {
            this.axios.defaults.headers.common.Authorization = token;
            this.updateUserInfo().then(() => {
                this.watch.ready = true;
            });
        }
    }

    public async whenBecomeReady(callback: () => void) {
        while (!this.ready()) {
            await this.sleep();
        }
        callback();
    }

    public ready(): boolean {
        return this.watch.ready;
    }

    public authenticated(): boolean {
        return this.watch.authenticated;
    }

    public user(): UserInfo {
        return this.watch.userInfo!;
    }

    public async login(username: string, password: string): Promise<Error | undefined> {
        try {
            const response = await this.axios.post('/auth/login', {
                username,
                password,
            });
            const token = response.headers.authorization;
            this.storeToken(token);
            this.axios.defaults.headers.common.Authorization = token;
            return await this.updateUserInfo();
        } catch (error) {
            return responseToError(error);
        }
    }

    public async signUp(
            firstName: string,
            lastName: string,
            username: string,
            password: string): Promise<Error | undefined> {
        try {
            const response = await this.axios.post('/auth/signup', {
                first_name: firstName,
                last_name: lastName,
                username,
                password,
            });
        } catch (error) {
            return responseToError(error);
        }
    }

    public logout(): void {
        this.watch.authenticated = false;
        this.watch.userInfo = undefined;
        this.removeToken();
        this.axios.defaults.headers.common.Authorization = undefined;
    }

    public async changePassword(oldPassword: string, newPassword: string): Promise<Error | undefined> {
        try {
            const response = await this.axios.post('/auth/change/password', {
                old_password: oldPassword,
                new_password: newPassword,
            });
        } catch (error) {
            return responseToError(error);
        }
    }

    private sleep(): Promise<void> {
        return new Promise((res) => setTimeout(res, 100 /*ms*/));
    }

    private async updateUserInfo(): Promise<Error | undefined> {
        try {
            const response = await this.axios.get('/auth/user');
            this.watch.userInfo = new UserInfo(response.data.data);
            this.watch.authenticated = true;
        } catch (error) {
            this.watch.authenticated = false;
            this.watch.userInfo = undefined;
            this.removeToken();
            this.axios.defaults.headers.common.Authorization = undefined;
            return responseToError(error);
        }
    }

    private storeToken(token: string): void {
        window.localStorage.setItem('auth_header_value', token);
    }

    private removeToken(): void {
        window.localStorage.removeItem('auth_header_value');
    }
}
