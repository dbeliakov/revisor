<template>
  <div style="margin: 20px 40px;">
    <div class="ui four item menu" style="float: right; width: 400px; font-size: 8pt;">
      <a class="item" :class="{active: showMode === 'all'}" @click.prevent="show('all')">Все</a>
      <a class="item" :class="{active: showMode === 'open'}" @click.prevent='show("open")'>Открытые</a>
      <a class="item" :class="{active: showMode === 'accepted'}" @click.prevent='show("accepted")'>Принятые</a>
      <a class="item" :class="{active: showMode === 'rejected'}" @click.prevent='show("rejected")'>Отклоненные</a>
    </div>
    <table class="ui single line striped table">
      <thead>
        <tr><th>
          Имя
        </th>
        <th>
          Создатель
        </th>
        <th>
          Ревьюеры
        </th>
        <th>
          Комментарии
        </th>
        <th>
          Обновлено
        </th>
      </tr></thead>
      <tbody>
        <tr v-for="review in computedFilteredReviews" :key="review.id"
            v-bind:class="{positive: review.closed && review.accepted, negative: review.closed && !review.accepted}">
          <td><router-link :to="'/review/' + review.id">{{review.name}}</router-link></td>
          <td class="collapsing">
            {{ review.owner.firstName }} {{ review.owner.lastName }}
          </td>
          <td class="collapsing">
            <span v-for="reviewer in review.reviewers" :key="reviewer.username">
              {{ reviewer.firstName }} {{ reviewer.lastName }}
            </span>
          </td>
          <td class="collapsing">
            <i class="comments outline icon"></i>{{review.commentsCount}}
            <i class="bug icon"></i>{{0}}
          </td>
          <td class="right aligned collapsing">{{timeToString(review.updated)}}</td>
        </tr>
      </tbody>
    </table>
  </div>
</template>

<script lang="ts">
import {Component, Vue, Prop, Watch} from 'vue-property-decorator';
import Review from '@/reviews/review';
import ReviewsService from '@/reviews/service';
import {timeToString} from '@/utils/utils';

@Component
export default class ReviewsList extends Vue {
    @Prop(String)
    public type!: string;

    public reviews: Review[] = [];
    public error: string = '';
    public showMode: string = 'open';

    public timeToString = timeToString;

    public created() {
        this.updateReviews();
    }

    public async updateReviews() {
        let result: Review[] | Error;
        if (this.type === 'incoming') {
            result = await this.$reviews.loadIncomingReviews();
        } else if (this.type === 'outgoing') {
            result = await this.$reviews.loadOutgoingReviews();
        } else {
            return;
        }
        if (result instanceof Error) {
            this.error = result.message;
            return;
        }
        this.reviews = result;
        this.error = '';
    }

    public get computedFilteredReviews() {
        if (this.showMode === 'all') {
            return this.reviews;
        } else if (this.showMode === 'open') {
          return this.reviews.filter((review) => !review.closed);
        } else if (this.showMode === 'accepted') {
          return this.reviews.filter((review) => review.closed && review.accepted);
        } else if (this.showMode === 'rejected') {
          return this.reviews.filter((review) => review.closed && !review.accepted);
        }
        return [];
    }

    public show(mode: string) {
      this.showMode = mode;
    }

    @Watch('$route')
    public onRouteChanged(from: any, to: any) {
        this.reviews = [];
        this.updateReviews();
    }
}
</script>

<style lang="scss">
@import '~semantic-ui-css/semantic.min.css';
</style>
