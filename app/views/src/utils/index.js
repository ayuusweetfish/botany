export default {
  install (Vue, options) {
    Vue.prototype.$functions = {
      dateStrings (timestamp) {
        let date = new Date(timestamp)
        let Y = date.getFullYear()
        let M = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1
        let D = date.getDate() < 10 ? '0' + date.getDate() : date.getDate()
        let h = date.getHours() < 10 ? '0' + date.getHours() : date.getHours()
        let m = date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()
        let s = date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds()
        return {
          'yyyy': Y,
          'MM': M,
          'dd': D,
          'HH': h,
          'mm': m,
          'ss': s
        }
      },
      dateTimeString (timestamp) {
        let date = new Date(timestamp)
        let Y = date.getFullYear()
        let M = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1
        let D = date.getDate() < 10 ? '0' + date.getDate() : date.getDate()
        let h = date.getHours() < 10 ? '0' + date.getHours() : date.getHours()
        let m = date.getMinutes() < 10 ? '0' + date.getMinutes() : date.getMinutes()
        let s = date.getSeconds() < 10 ? '0' + date.getSeconds() : date.getSeconds()
        return Y + '-' + M + '-' + D + ' ' + h + ':' + m + ':' + s
      },
      dateString (timestamp) {
        let date = new Date(timestamp)
        let Y = date.getFullYear()
        let M = date.getMonth() + 1 < 10 ? '0' + (date.getMonth() + 1) : date.getMonth() + 1
        let D = date.getDate() < 10 ? '0' + date.getDate() : date.getDate()
        return Y + '-' + M + '-' + D
      },
      globalValidator (rule, value, callback) {
        if (value.match(new RegExp('[~#^$%&!?%* ]'))) {
          callback(new Error('输入了非法字符'))
        } else {
          callback()
        }
      }
    }
    Vue.prototype.$consts = {
      privilege: {
        'common': 0,
        'organizer': 1,
        'superuser': 2
      },
      role: {
        'notIn': -1,
        'imIn': 0,
        'moderator': 1
      }
    }
  }
}
