<template>
  <div class="help-page">
    <a-card>
      <a-typography>
        <a-typography-title :level="3">淮北师范大学学生助手</a-typography-title>
        <a-typography-paragraph type="secondary">
          一款为淮北师范大学同学打造的桌面小工具，帮你快速查询常用校园信息。
        </a-typography-paragraph>

        <a-typography-title :level="4">功能介绍</a-typography-title>
        <a-typography-paragraph>
          <ul>
            <li><strong>成绩查询</strong> — 登录教务系统，按学期查看各科成绩和 GPA</li>
            <li><strong>课表查询</strong> — 查看每学期的课程安排</li>
            <li><strong>校园新闻</strong> — 浏览学校最新资讯</li>
            <li><strong>第二课堂成绩</strong> — 查看五育模块（思想政治、专业技能、体育健身、文化艺术、志愿服务）学时统计</li>
          </ul>
        </a-typography-paragraph>

        <a-typography-title :level="4">使用流程</a-typography-title>
        <a-typography-paragraph>
          <ol>
            <li>前往侧边栏的「<strong>设置</strong>」页面，配置你的账号密码</li>
            <li>教务系统和校园门户是<strong>两套独立账号</strong>，需要分别配置</li>
            <li>配置完成后，点击对应功能页面即可自动登录并查询</li>
          </ol>
        </a-typography-paragraph>

        <a-divider />

        <a-typography-paragraph type="secondary">
          这个 App 还在持续开发中，后续会加入更多实用功能。如果你有什么好想法或建议，欢迎联系我：
          <a href="mailto:zguigo@163.com">zguigo@163.com</a>
        </a-typography-paragraph>
      </a-typography>
    </a-card>
  </div>
</template>

<style scoped>
.help-page {
  max-width: 700px;
}

.help-page ul,
.help-page ol {
  padding-left: 20px;
}

.help-page li {
  margin-bottom: 8px;
  line-height: 1.8;
}
</style>
