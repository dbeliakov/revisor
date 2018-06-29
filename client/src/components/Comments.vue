<template>
    <div>
        <div v-for="comment in comments" :key="comment.id" class="comment">
            <div class="header"><span class="login">{{ comment.author.first_name }} {{ comment.author.last_name }} ({{ comment.author.username }})</span><span class="date">{{ timeConverter(comment.created) }}</span></div>
            <div class="text">{{ comment.text }}</div>
            <div class="footer"><a href="#" @click.prevent="$emit('new-comment')">Ответить</a></div>
        </div>
    </div>
</template>

<script>
export default {
  name: 'Comments',
  props: ['comments'],
  methods: {
    timeConverter (timestamp) {
      var toStr = function (val) {
        if (val < 10) {
          return '0' + val
        }
        return '' + val
      }

      var a = new Date(timestamp * 1000)
      var months = ['Января', 'Февраля', 'Марта', 'Апреля', 'Мая', 'Июня', 'Июля', 'Августа',
        'Сентября', 'Октября', 'Ноября', 'Декабря']
      var year = a.getFullYear()
      var month = months[a.getMonth()]
      var date = a.getDate()
      var hour = a.getHours()
      var min = a.getMinutes()
      var time = date + ' ' + month + ' ' + year + ' ' + toStr(hour) + ':' + toStr(min)
      return time
    }
  }
}
</script>

<style scoped>
.login {
    font-weight: bold;
    font-size: 12pt;
}

.date {
    color: grey;
    margin-left: 20px;
}

.comment {
    padding: 20px;
    padding-top: 10px;
    padding-bottom: 10px;
    border: 1px solid grey;
    border-radius: 5px;
    margin-left: 64px;
}

.comment2 {
    padding: 20px;
    padding-top: 10px;
    padding-bottom: 10px;
    border: 1px solid grey;
    border-top: 0px solid black;
    border-radius: 5px;
    border-top-left-radius: 0px;
    border-top-right-radius: 0px;
    margin-left: 84px;
}

.header {
    margin-bottom: 10px;
}

.footer {
    margin-top: 10px;
}
</style>
