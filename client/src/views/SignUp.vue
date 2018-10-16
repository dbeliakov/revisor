<template>
    <div id="login">
        <div class="column">
            <form class="ui large form" v-on:submit.prevent="signUp()">
                <div class="ui">
                <div class="field">
                    <div class="ui input">
                    <input name="firstName" v-model="firstName" placeholder="Имя" type="text">
                    </div>
                </div>
                <div class="field">
                    <div class="ui input">
                    <input name="lastName" v-model="lastName" placeholder="Фамилия" type="text">
                    </div>
                </div>
                <div class="field">
                    <div class="ui input">
                    <input name="username" v-model="username" placeholder="Логин" type="text">
                    </div>
                </div>
                <div class="field">
                    <div class="ui input">
                    <input name="password" v-model="password" placeholder="Пароль" type="password">
                    </div>
                </div>
                <button class="ui fluid large blue submit button" v-bind:class="{'disabled': formDisabled}" type="submit">Зарегистрироваться</button>
                </div>

                <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>

            </form>
            <div class="ui message">
                <router-link :to="{name: 'login'}">Вход</router-link>
            </div>
        </div>
    </div>
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator';

@Component
export default class SignUp extends Vue {
    public firstName: string = '';
    public lastName: string = '';
    public username: string = '';
    public password: string = '';
    public error: string = '';
    public formDisabled: boolean = false;

    public async signUp(): Promise<void> {
        try {
            this.disableForm();
            if (!this.validateForm()) {
                return;
            }
            const error = await this.$auth.signUp(
                this.firstName,
                this.lastName,
                this.username,
                this.password);
            if (error) {
                this.error = error.message;
            } else {
                this.$router.push({name: 'login'});
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
        if (this.firstName.length === 0) {
            this.error = 'Введите имя';
            return false;
        }
        if (this.lastName.length === 0) {
            this.error = 'Введите фамилию';
        }
        if (this.username.length === 0) {
            this.error = 'Введите логин';
            return false;
        }
        if (this.password.length < 6) {
            this.error = 'Пароль должен быть не короче 6 символов';
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
