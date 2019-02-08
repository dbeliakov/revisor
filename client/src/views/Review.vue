<template>
  <div style="margin: 20px 40px; text-align: left;">
    <div v-if="data">
      <h1>{{ data.info.name }}</h1>
      <h4 style="display: inline;">Создатель:</h4> {{ data.info.owner.firstName }}
        {{ data.info.owner.lastName }} ({{ data.info.owner.username }})<br>
      <h4 style="display: inline;">Ревьюеры:</h4>
        <span v-for="reviewer in data.info.reviewers" :key="reviewer.username">
            {{reviewer.firstName}} {{reviewer.lastName}} ({{reviewer.username}})
            <template v-if="reviewer !== data.info.reviewers[data.info.reviewers.length - 1]">,</template>
        </span><br>
      <template v-if="data.info.closed"><h4 style="display: inline;">Закрыто:</h4> <span v-if="data.info.accepted"> Принято</span> <span v-if="!data.info.accepted"> Отклонено</span><br></template>
      <span><h4 style="display: inline;">Обновлено:</h4> {{timeToString(data.info.updated)}}</span>
      <div v-if="!data.info.closed" style="margin-top: 20px;">
        <button v-if="data.info.owner.username === $auth.user().username" class="review-button ui primary basic button" @click="openModal">Обновить</button>
        <button v-if="data.info.owner.username !== $auth.user().username" class="review-button ui positive basic button" @click="accept">Принять</button>
        <button class="review-button ui negative basic button" @click="decline">Отклонить</button><br>
      </div>
    </div>
    <div v-if="data" style="text-align: center; margin-bottom: 20px;">
      <h4 style="margin-top: 30px; margin-bottom: 10px;" v-if="data.info.revisionsCount > 1 && startRev !== null && endRev !== null">
        <span v-if="startRev != endRev">revision {{startRev}} <i class="right arrow icon"></i> revision {{endRev}}</span>
        <span v-else>revision {{startRev}}</span>
      </h4>
      <div v-if="data.info.revisionsCount > 1" id="slider_wrapper" style="margin: 0px auto; width:300px;">
        <div id="slider_revisions"></div>
      </div>
    </div>

    <DiffComponent v-if="data" :diff="data.diff" :commentsList="data.comments" :reviewId="$route.params.id" @update-all="loadData"></DiffComponent>

    <div class="ui modal" id="add_revision">
      <i class="close icon"></i>
      <div class="header">
        Обновить ревью
      </div>
      <div class="content">
        <form class="ui large form" v-on:submit.prevent="updateReview()">
        <div class="ui">
          <div class="field">
            <div class="ui input">
              <input name="name" v-model="name" placeholder="Заголовок" type="text">
            </div>
          </div>
          <div class="reviewers-dropdown">
            <div class="field" style="margin-bottom: 0px !important;">
                <div class="ui input">
                    <input
                        name="reviewers"
                        v-model="reviewerSearch"
                        placeholder="Добавить ревьюера"
                        type="text"
                        autocomplete="off"
                        @input="loadReviewersList"
                        @keydown.enter.prevent>
                </div>
            </div>
            <div class="ui vertical menu reviewers-list" style="width: 100%;" v-if="reviewersList.length > 0">
                <a class="item"
                    v-for="reviewer in reviewersList"
                    :key="reviewer.id"
                    @click="addReviewer(reviewer.username)">
                    {{reviewer.firstName}} {{reviewer.lastName}} ({{reviewer.username}})
                </a>
            </div>
          </div>
          <div class="field" style="text-align: left; font-size: 14pt; margin-top: 20px;">
                <b>Ревьюеры: </b><span v-if="reviewers.length === 0">
                  нет добавленных ревьюеров
                </span><span v-for="reviewer in reviewers" :key="reviewer.id">{{reviewer}}
                    <i style="font-size: 12pt; cursor: pointer;"
                        @click="removeReviewer(reviewer)"
                        class="window close outline icon"></i></span>
          </div>
          <FileLoader
            @onStartReading='onStartReadingFile'
            @onFinishReading='onFinishReadingFile'
            @onReadingError='onReadingError'/>
          <button class="ui fluid large blue submit button" v-bind:class="{'disabled': formDisabled}" type="submit">Обновить ревью</button>
        </div>
        <div class="ui negative message" v-if="error.length > 0">{{ error }}</div>
      </form>
      </div>
    </div>
  </div>
</template>

<script lang="ts">
import {Component, Vue, Watch} from 'vue-property-decorator';
import DiffComponent from '@/components/Diff.vue';
import FileLoader from '@/components/FileLoader.vue';
import { DiffReply } from '@/reviews/service';
import { Diff } from '@/reviews/diff';
import {timeToString} from '@/utils/utils';

import * as $ from 'jquery';
import { UserInfo } from '../auth/user-info';
(window as any).$ = $;
(window as any).jQuery = $;

/* tslint:disable:no-var-requires */
require('semantic-ui-css/semantic.min.js');
require('jquery-ui/ui/widgets/slider.js');
require('jquery-ui/themes/base/all.css');

@Component({
  components: {DiffComponent, FileLoader},
})
export default class Review extends Vue {
  // public reviewId: string = '';
  public data: DiffReply | null = null;
  public startRev: number | null = null;
  public endRev: number | null = null;

  public name: string | null = null;
  public filename: string | null = null;
  public reviewerSearch: string = '';
  public reviewersList: UserInfo[] = [];
  public reviewers: string[] = [];
  public formDisabled: boolean = false;
  public fileContent: string = '';
  public error: string = '';

  public timeToString = timeToString;

  public created() {
    this.loadData();
  }

  @Watch('$route')
  public onChildChanged(to: string, from: string) {
    this.data = null;
    this.loadData();
  }

  public async loadData() {
    const result = await this.$reviews.loadDiff(+this.$route.params.id, this.startRev, this.endRev);
    if (result instanceof Error) {
      alert(result);
      return;
    }
    this.data = result;
    if (!this.startRev) {
      this.startRev = 1;
    }
    if (!this.endRev) {
      this.endRev = this.data.info.revisionsCount;
    }
    const that = this;
    if (this.data.info.revisionsCount > 1) {
      $(() => {
        // TODO no slider without any
        ($('#slider_revisions') as any).slider({
          range: true,
          min: 1,
          max: that.data!.info.revisionsCount,
          values: [that.startRev, that.endRev],
          // TODO add types
          slide(event: any, ui: any) {
            that.startRev = ui.values[0];
            that.endRev = ui.values[1];
            that.loadData();
          },
        });
        const width = Math.min(800, this.data!.info.revisionsCount * 40);
        $('#slider_wrapper').css('width', '' + width + 'px');
      });
    }
  }

  public openModal() {
    // TODO no modal without any
    this.name = this.data!.info.name;
    this.reviewers = this.data!.info.reviewers.map((reviewer: UserInfo) => reviewer.username);
    this.filename = this.data!.diff.filename;
    ($('#add_revision') as any).modal('show');
  }

  public onStartReadingFile() {
    this.disableForm();
  }

  public onFinishReadingFile(filename: string, content: string) {
    this.fileContent = content.replace(/^data:.+;base64,/, '');
    this.filename = filename;
    this.error = '';
    this.enableForm();
  }

  public onReadingError(err: string) {
    this.filename = '';
    this.fileContent = '';
    this.enableForm();
    alert(err);
  }

  public addReviewer(username: string) {
    if (this.reviewers.indexOf(username) === -1) {
      this.reviewers.push(username);
    }
    this.reviewersList = [];
    this.reviewerSearch = '';
  }

  public removeReviewer(username: string) {
    const index = this.reviewers.indexOf(username);
    if (index > -1) {
      this.reviewers.splice(index, 1);
    }
  }

  public async loadReviewersList() {
    if (this.reviewerSearch.length === 0) {
      this.reviewersList = [];
      return;
    }
    const result = await this.$reviews.searchReviewers(this.reviewerSearch);
    if (result instanceof Error) {
      alert(result);
      return;
    }
    this.reviewersList = result;
  }

  public async updateReview() {
    if (this.formDisabled) {
      return;
    }
    this.disableForm();
    this.error = '';
    try {
      if (!this.validateForm()) {
        return;
      }
      const result = await this.$reviews.updateReview(
        +this.$route.params.id, this.name!, this.reviewers.join(','), this.filename, this.fileContent);
      if (result) {
        this.error = result.message;
      } else {
        this.filename = null;
        this.fileContent = '';
        // TODO no modal without any
        ($('#add_revision') as any).modal('hide');
        this.loadData();
      }
    } finally {
      this.enableForm();
    }
  }

  public async accept() {
    const result = await this.$reviews.acceptReview(+this.$route.params.id);
    if (result) {
      alert(result);
    }
    this.loadData();
  }

  public async decline() {
    const result = await this.$reviews.declineReview(+this.$route.params.id);
    if (result) {
      alert(result);
    }
    this.loadData();
  }

  private validateForm(): boolean {
    if (this.name!.length === 0) {
      this.error = 'Введите заголовок ревью';
      return false;
    }
    if (this.reviewers.length === 0) {
      this.error = 'Добавьте как минимум одного ревьюера';
      return false;
    }
    return true;
  }

  private disableForm() {
    this.formDisabled = true;
  }

  private enableForm() {
    this.formDisabled = false;
  }
}
</script>