webpackJsonp([42606486476073],{345:function(e,t){e.exports={data:{markdownRemark:{html:'<p>In a production deployment you\'ll want multiple instances of the flotilla service running and postgres running elsewhere (eg. Amazon RDS). In this case the most salient detail configuration detail is the <code>DATABASE_URL</code>.</p>\n<h3>Docker based deploy</h3>\n<p>The simplest way to deploy for very light usage is to avoid a reverse proxy and deploy directly with docker.</p>\n<ol>\n<li>\n<p>Build and tag an image for flotilla using the <code>Dockerfile</code> provided in this repo:</p>\n<pre><code>docker build -t &#x3C;your repo name>/flotilla:&#x3C;version tag>\n</code></pre>\n</li>\n<li>\n<p>Run this image wherever you deploy your services:</p>\n<pre><code>docker run -e DATABASE_URL=&#x3C;your db url> -e FLOTILLA_MODE=prod -p 3000:3000 ...&#x3C;other standard docker run args>\n</code></pre>\n<blockquote>\n<h2>Notes:</h2>\n<ul>\n<li>Flotilla uses <a href="https://github.com/spf13/viper">viper</a> for configuration so you can override any of the default configuration under <code>conf/</code> using run time environment variables passed to <code>docker run</code></li>\n<li>In most realistic deploys you\'ll likely want to configure a reverse proxy to sit in front of the flotilla container. See the docs <a href="https://hub.docker.com/_/nginx/">here</a></li>\n</ul>\n</blockquote>\n</li>\n</ol>\n<pre><code>See [docker run](https://docs.docker.com/engine/reference/run/) for more details\n</code></pre>\n<h3>Configuration In Detail</h3>\n<p>The variables in <code>conf/config.yml</code> are sensible defaults. Most should be left alone unless you\'re developing flotilla itself. However, there are a few you may want to change in a production environment.</p>\n<table>\n<thead>\n<tr>\n<th>Variable Name</th>\n<th>Description</th>\n</tr>\n</thead>\n<tbody>\n<tr>\n<td><code>worker.retry_interval</code></td>\n<td>Run frequency of the retry worker</td>\n</tr>\n<tr>\n<td><code>worker.submit_interval</code></td>\n<td>Poll frequency of the submit worker</td>\n</tr>\n<tr>\n<td><code>worker.status_interval</code></td>\n<td>Poll frequency of the status update worker</td>\n</tr>\n<tr>\n<td><code>http.server.read_timeout_seconds</code></td>\n<td>Sets read timeout in seconds for the http server</td>\n</tr>\n<tr>\n<td><code>http.server.write_timeout_seconds</code></td>\n<td>Sets the write timeout in seconds for the http server</td>\n</tr>\n<tr>\n<td><code>http.server.listen_address</code></td>\n<td>The port for the http server to listen on</td>\n</tr>\n<tr>\n<td><code>owner_id_var</code></td>\n<td>Which environment variable containing ownership information to inject into the runtime of jobs</td>\n</tr>\n<tr>\n<td><code>enabled_workers</code></td>\n<td>This variable is a list of the workers that run. Use this to control what workers run when using a multi-container deployment strategy. Valid list items include (\n<code>retry</code>\n, \n<code>submit</code>\n, and \n<code>status</code>\n)</td>\n</tr>\n<tr>\n<td><code>log.namespace</code></td>\n<td>For the default ECS execution engine setup this is the \n<code>log-group</code>\n to use</td>\n</tr>\n<tr>\n<td><code>log.retention_days</code></td>\n<td>For the default ECS execution engine this is the number of days to retain logs</td>\n</tr>\n<tr>\n<td><code>log.driver.options.*</code></td>\n<td>For the default ECS execution engine these map to the \n<code>awslogs</code>\n driver options \n<a href="https://docs.aws.amazon.com/AmazonECS/latest/developerguide/using_awslogs.html">here</a></td>\n</tr>\n<tr>\n<td><code>queue.namespace</code></td>\n<td>For the default ECS execution engine this is the prefix used for SQS to determine which queues to pull job launch messages from</td>\n</tr>\n<tr>\n<td><code>queue.retention_seconds</code></td>\n<td>For the default ECS execution engine this configures how long a message will stay in an SQS queue without being consumed</td>\n</tr>\n<tr>\n<td><code>queue.process_time</code></td>\n<td>For the default ECS execution engine configures the length of time allowed to process a job launch message</td>\n</tr>\n<tr>\n<td><code>queue.status</code></td>\n<td>For the default ECS execution engine this configures which SQS queue to route ECS cluster status updates to</td>\n</tr>\n<tr>\n<td><code>queue.status_rule</code></td>\n<td>For the default ECS execution engine this configures the name of the rule for routing ECS cluster status updates</td>\n</tr>\n</tbody>\n</table>',frontmatter:{path:"/docs/deploying",title:"Deploying",group:"docs",index:"4"}}},pathContext:{}}}});
//# sourceMappingURL=path---docs-deploying-577a7e3c520f4e2009e5.js.map