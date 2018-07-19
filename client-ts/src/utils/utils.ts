export function responseToError(error: any): Error {
    if (!error.response || !error.response.status) {
        return new Error('Ошибка сети');
    } else if (error.response.data.client_message) {
        return new Error(error.response.data.client_message);
    } else {
        return new Error('Внутренняя ошибка сервиса');
    }
}

export function timeToString(time: Date): string {
    const numberToString = (value: number) => {
      if (value < 10) {
        return '0' + value.toString();
      }
      return value.toString();
    };

    const months = ['Января', 'Февраля', 'Марта', 'Апреля', 'Мая', 'Июня', 'Июля', 'Августа',
      'Сентября', 'Октября', 'Ноября', 'Декабря'];
    const year = time.getFullYear();
    const month = months[time.getMonth()];
    const day = time.getDate();
    const hour = time.getHours();
    const min = time.getMinutes();
    return day.toString() + ' ' + month + ' ' + year.toString() + ' ' +
        numberToString(hour) + ':' + numberToString(min);
}
