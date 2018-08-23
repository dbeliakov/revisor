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
  </div>
</template>

<script lang="ts">
import {Component, Vue} from 'vue-property-decorator';
import { error } from 'util';

@Component
export default class Profile extends Vue {
  public oldPassword: string = '';
  public newPassword: string = '';
  public error: string = '';
  public formDisabled: boolean = false;

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
