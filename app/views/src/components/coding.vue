<template>
  <div>
    <el-row style="margin-bottom: 10px">
      <el-card>
        <div align="left">
          <div style="display: inline">代码已提交，状态：</div>
          <div style="display: inline; color: orange">处理中</div>
        </div>
      </el-card>
    </el-row>
    <el-row :gutter="20" style="margin-bottom: 10px">
      <el-col :span="18">
        <el-card body-style="height: 360px">
          <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">代码编辑</div>
          <codemirror
            v-model="code"
            :options="cmOptions"
            class="code"
            align="left"
          >
          </codemirror>
        </el-card>
      </el-col>
      <el-col :span="6">
        <el-card body-style="height: 360px">
          <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">历史代码</div>
          <el-timeline align="left">
            <el-timeline-item
              v-for="(activity, index) in history"
              :key="index"
              :timestamp="activity.time"
              :color="activity.color"
              placement="top"
            >
            <div align="left" style="display: inline">{{activity.status}}</div>
            <el-button type="text" size="small" style="display: inline">点击导出</el-button>
            </el-timeline-item>
          </el-timeline>
        </el-card>
      </el-col>
    </el-row>
    <el-row>
      <el-card>
        <div align="left" style="font-size: 18px; font-weight: 600; margin-bottom: 20px">操作</div>
        <el-row>
          <el-col :span="8">
            <el-button type="primary" size="small" style="width: 80%">提交</el-button>
          </el-col>
          <el-col :span="8">
            <el-button size="small" style="width: 80%">保存至草稿</el-button>
          </el-col>
          <el-col :span="8">
            <el-button size="small" style="width: 80%">导出草稿</el-button>
          </el-col>
        </el-row>
      </el-card>
    </el-row>
  </div>
</template>

<script>
import {codemirror} from 'vue-codemirror-lite'

export default {
  name: 'coding',
  components: {
    codemirror
  },
  data () {
    return {
      code: '',
      history: [
        {
          time: '2019-10-25 11:00:00',
          status: '处理中',
          color: 'orange'
        },
        {
          time: '2019-09-10 19:02:03',
          status: '编译失败',
          color: 'red'
        },
        {
          time: '2019-09-10 19:00:00',
          status: '可用',
          color: 'green'
        }
      ],
      cmOptions: {
        lineNumbers: true,
        indentUnit: 2,
        autoCloseBrackets: true
      }
    }
  }
}
</script>

<style scoped>
  .code{
    border: 1px solid #dcdfe6;
    margin: auto;
  }
</style>
