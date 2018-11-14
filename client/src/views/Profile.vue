<template>
  <div id="profile">
    <h1>Смена пароля</h1>
    <div class="column">
      <form class="ui large form" v-on:submit.prevent="changePassword()">
        <div class="ui">
          <div class="field">
            <div class="ui input">
              <input name="password" v-model="oldPassword" placeholder="Текущий пароль" type="password">
            </div>
          </div>
          <div class="field">
            <div class="ui input">
              <input name="password" v-model="newPassword" placeholder="Новый пароль" type="password">
            </div>
          </div>
          <button class="ui fluid large blue submit button" :class="{'disabled': formDisabled}" type="submit">Сменить пароль</button>
        </div>
        <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>
      </form>
    </div>
    
    <h1 style="margin-top: 60px;">Уведомления в Telegram</h1>
    <span style="font-size: 15pt;" v-if="tgLoaded && tgUsername != null">
      Подключен аккаунт <a :href='"https://t.me/" + $auth.user().tgUsername'>@{{$auth.user().tgUsername}}</a>.
      <button class="ui small blue button" @click="unlinkTelegram">Отключить</button>
    </span>
    <div v-if="tgLoaded && !tgUsername" id="telegram"></div>
  </div>
</template>

<script lang="ts">
import {Component, Vue, Watch} from 'vue-property-decorator';
import { error } from 'util';

@Component
export default class Profile extends Vue {
  public oldPassword: string = '';
  public newPassword: string = '';
  public error: string = '';
  public formDisabled: boolean = false;
  public tgUsername: string | null = null;
  public tgLoaded: boolean = false;

  public async changePassword() {
    if (this.formDisabled) {
      return;
    }
    this.disableForm();
    try {
      if (!this.validateForm()) {
        return;
      }
      const result = await this.$auth.changePassword(this.oldPassword, this.newPassword);
      if (result) {
        this.error = result.message;
        return;
      }
      this.$router.push({name: 'home'});
    } finally {
      this.enableForm();
    }
  }

  public mounted() {
    this.updateTgButton();
  }

  public async updateTgButton() {
    this.tgLoaded = false;
    const login = await this.$auth.telegramLogin();
    if (login instanceof Error) {
      alert(error);
      return;
    }
    if (login != null) {
      this.tgUsername = login;
    } else {
      this.tgUsername = null;
    }
    this.tgLoaded = true;
  }

  public async unlinkTelegram() {
    const response = await this.$auth.unlinkTelegram();
    if (response) {
      alert(response);
    } else {
      this.updateTgButton();
    }
  }

  public updated() {
    if (this.tgUsername != null) {
      return;
    }
    const script = document.createElement('script');
    script.async = true;
    script.src = 'https://telegram.org/js/telegram-widget.js?5';
    script.setAttribute('data-size', 'large');
    script.setAttribute('data-telegram-login', 'RevisorNotificationsBot');
    script.setAttribute('data-request-access', 'write');
    script.setAttribute('data-onauth', 'onTelegramAuth(user)');
    (window as any).onTelegramAuth = this.onTelegramAuth;
    const target = document.getElementById('telegram');
    if (target == null) {
      return;
    }
    target!.appendChild(script);
  }

  private async onTelegramAuth(user: any) {
    const response = await this.$auth.linkTelegram(user.username, user.id);
    if (response) {
      alert(response);
    }
    this.updateTgButton();
  }

  private disableForm(): void {
    this.formDisabled = true;
  }

  private enableForm(): void {
    this.formDisabled = false;
  }

  private validateForm(): boolean {
    if (this.oldPassword.length < 6) {
      this.error = 'Введите текущий пароль (не короче 6 символов)';
      return false;
    }
    if (this.newPassword.length < 0) {
      this.error = 'Введите новый пароль (не короче 6 символов)';
      return false;
    }
    return true;
  }
}
</script>
