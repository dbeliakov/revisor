<template>
    <div id="login">
        <div class="column">
            <form class="ui large form" v-on:submit.prevent="login()">
                <div class="ui">
                <div class="field">
                    <div class="ui left icon input">
                    <i class="user icon"></i>
                    <input name="username" v-model="username" placeholder="Логин" type="text">
                    </div>
                </div>
                <div class="field">
                    <div class="ui left icon input">
                    <i class="lock icon"></i>
                    <input name="password" v-model="password" placeholder="Пароль" type="password">
                    </div>
                </div>
                <button class="ui fluid large blue submit button" v-bind:class="{'disabled': formDisabled}" type="submit">Войти</button>
                </div>

                <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>

            </form>
            <div class="ui message">
                <router-link :to="{name: 'signup'}">Регистрация</router-link>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator';

@Component
export default class Login extends Vue {
    public username: string = '';
    public password: string = '';
    public error: string = '';
    public formDisabled: boolean = false;

    public async login(): Promise<void> {
        try {
            this.disableForm();
            if (!this.validateForm()) {
                return;
            }
            const error = await this.$auth.login(this.username, this.password);
            if (error) {
                this.error = error.message;
            } else {
                this.$router.push({name: 'home'});
            }
        } finally {
            this.enableForm();
        }
    }

    private disableForm(): void {
        this.formDisabled = true;
    }

    private enableForm(): void {
        this.formDisabled = false;
    }

    private validateForm(): boolean {
        if (this.username.length === 0) {
            this.error = 'Введите логин';
            return false;
        }
        if (this.password.length === 0) {
            this.error = 'Введите пароль';
            return false;
        }
        return true;
    }
}
</script>

<style scoped lang="scss">
@import '~semantic-ui-css/semantic.min.css';

body > .grid {
    height: 100%;
}
.column {
    max-width: 450px;
    margin: auto;
    margin-top: 40px;
}
</style>
