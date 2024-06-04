# 分布式 ID 生成

1. 多主节点冗余，利用数据库的 auto_increment，但不是 +1，而是 +n，n 是主节点数量，这样就可以避免生成冲突的节点。缺点：不随时间递增；不能很好应对伸缩的情况
2. UUID，128 位的数字和字母，生成简单，不需要考虑协调和同步，易伸缩。缺点：跟时间没有关系；无法比较
3. Flicker’s ticket server，有一个中心化的数据库服务器，优势是纯数字，易实现，缺点是中心化导致的 single point of failure
4. Twitter’s Snowflake，综合了所有优点，生成的 ID 是一个 64 位的整数，主要有以下部分：timestamp，最重要的 41 位；datacenter ID, machine ID 都是一开始生成的；sequence ID，每毫秒从 0 开始，最多 4096 个

进阶问题：时钟同步（NTP）
