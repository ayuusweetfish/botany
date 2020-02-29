<template>
  <div class="d-flex mb-6">
    <v-icon class="mb-2">mdi-account-cog-outline</v-icon>
    <div class="d-flex flex-wrap">
      <div></div>
      <v-chip large outlined label close class="ml-2 mb-2"
        v-for="(item, index) in list" :key="index"
        @click:close="removeModerator(item.id)"
      >
        <user-tag size="small" :user="item" class="mr-2" disabled></user-tag>
      </v-chip>
      <v-menu
        v-model="selecting"
        offset-y
        :close-on-content-click="false"
      >
        <template v-slot:activator="{ on }">
          <v-chip large outlined label
            @click="selecting=true"
            v-on="on"
            color="primary"
            class="ml-2 mb-2"
            :disabled="list.length>20"
          ><v-icon class="body-1">mdi-plus</v-icon>添加管理员
        </v-chip>
        </template>
        <v-card :loading="searchLoading">
          <v-card-text>
            <div>
              <v-text-field
                outlined dense
                label="输入账户名搜索用户"
                v-model="searchText"
                append-outer-icon="mdi-magnify"
                :disabled="list.length>20||searchLoading"
                @click:append-outer="searchUsers"
                @keypress.enter="searchUsers"
              ></v-text-field>
            </div>
            <div
              v-if="selections.length===0"
              class="d-flex justify-center"
            >暂无数据</div>
            <div>
              <div
                v-for="(item, index) in selections" :key="index"
                class="d-flex justify-space-between align-center"
              >
                <user-tag size="small" :user="item" identify></user-tag>
                <v-btn
                  text
                  color="primary"
                  :disabled="list.length>20||isSelectionDisabled(item.id)"
                  @click="addModerator(item)"
                ><v-icon>mdi-plus-box</v-icon></v-btn>
              </div>
            </div>
          </v-card-text>
          <v-divider></v-divider>
          <v-card-actions class="justify-end">
            <v-btn text color="primary" @click="selecting=false">完成</v-btn>
          </v-card-actions>
        </v-card>
      </v-menu>
    </div>
  </div>
</template>

<script>
import UserTag from '../components/UserTag.vue'
export default {
  props: {
    value: Array
  },
  components: {
    'user-tag': UserTag
  },
  watch: {
    value: function (newval, oldval) {
      this.list = newval
    }
  },
  data: () => ({
    list: [],
    selections: [],
    selecting: false,
    searchText: '',
    searchLoading: false
  }),
  methods: {
    searchUsers () {
      if (this.searchText !== '') {
        this.searchLoading = true
        this.$axios.get(
          '/user_search/' + this.searchText
        ).then(res => {
          this.selections = res.data
          this.searchLoading = false
        }).catch(() => {
          this.searchLoading = false
        })
      } else {
        this.selections = []
      }
    },
    isSelectionDisabled (id) {
      if (id === this.$store.state.id) {
        return true
      }
      for (let i = 0; i < this.list.length; ++i) {
        if (this.list[i].id === id) {
          return true
        }
      }
      return false
    },
    removeModerator (id) {
      this.list.splice(this.list.findIndex(item => item.id === id), 1)
      this.emit()
    },
    addModerator (moderator) {
      this.list.push(moderator)
      this.emit()
    },
    emit () {
      this.$emit('input', this.list)
      this.$emit('change')
    }
  }
}
</script>

<style>

</style>
