<?php

declare(strict_types=1);
/**
 * This file is part of MineAdmin.
 *
 * @link     https://www.mineadmin.com
 * @document https://doc.mineadmin.com
 * @contact  root@imoi.cn
 * @license  https://github.com/mineadmin/MineAdmin/blob/master/LICENSE
 */

namespace Mine\Command;

use Hyperf\Command\Annotation\Command;
use Hyperf\Database\Schema\Schema;
use Hyperf\DbConnection\Db;
use Mine\Mine;
use Mine\MineCommand;

use function Hyperf\Support\env;
use function Hyperf\Support\make;

/**
 * Class InstallProjectCommand.
 */
#[Command]
class InstallProjectCommand extends MineCommand
{
    protected const CONSOLE_GREEN_BEGIN = "\033[32;5;1m";

    protected const CONSOLE_RED_BEGIN = "\033[31;5;1m";

    protected const CONSOLE_END = "\033[0m";

    /**
     * 安装命令.
     */
    protected ?string $name = 'mine:install';

    protected array $database = [];

    protected array $redis = [];

    public function configure()
    {
        parent::configure();
        $this->setHelp('run "php bin/hyperf.php mine:install" install MineAdmin system');
        $this->setDescription('MineAdmin system install command');
    }

    public function handle(): void
    {
        $this->installLocalModule();
        $this->setOthers();
        $this->finish();
    }

    protected function welcome(): void
    {
        $this->line('-----------------------------------------------------------', 'comment');
        $this->line('Hello, welcome use MineAdmin system.', 'comment');
        $this->line('The installation is about to start, just a few steps', 'comment');
        $this->line('-----------------------------------------------------------', 'comment');
    }

    /**
     * install modules.
     */
    protected function installLocalModule(): void
    {
        /* @var Mine $mine */
        $this->line("Installation of local modules is about to begin...\n", 'comment');
        $mine = make(Mine::class);
        $modules = $mine->getModuleInfo();
        foreach ($modules as $name => $info) {
            $this->call('mine:migrate-run', ['name' => $name, '--force' => 'true']);
            if ($name === 'System') {
                $this->initUserData();
            }
            $this->call('mine:seeder-run', ['name' => $name, '--force' => 'true']);
            $this->line($this->getGreenText(sprintf('"%s" module install successfully', $name)));
        }
    }

    protected function setOthers(): void
    {
        $this->line(PHP_EOL . ' MineAdmin set others items...' . PHP_EOL, 'comment');
        $this->call('mine:update');
        $this->call('mine:jwt-gen', ['--jwtSecret' => 'JWT_SECRET']);
        $this->call('mine:jwt-gen', ['--jwtSecret' => 'JWT_API_SECRET']);

        if (! file_exists(BASE_PATH . '/config/autoload/mineadmin.php')) {
            $this->call('vendor:publish', ['package' => 'xmo/mine']);
        }

        $downloadFrontCode = $this->confirm('Do you downloading the front-end code to "./web" directory?', true);

        // 下载前端代码
        if ($downloadFrontCode) {
            $this->line(PHP_EOL . ' Now about to start downloading the front-end code' . PHP_EOL, 'comment');
            if (\shell_exec('which git')) {
                \system('git clone https://gitee.com/mineadmin/mineadmin-vue.git ./web/');
            } else {
                $this->warn('Your server does not have the `git` command installed and will skip downloading the front-end project');
            }
        }
    }

    protected function initUserData(): void
    {
        // 清理数据
        Db::table('system_user')->truncate();
        Db::table('system_role')->truncate();
        Db::table('system_user_role')->truncate();
        if (Schema::hasTable('system_user_dept')) {
            Db::table('system_user_dept')->truncate();
        }

        // 创建超级管理员
        $superAdminId = Db::table('system_user')->insertGetId([
            'username' => 'superAdmin',
            'password' => password_hash('admin123', PASSWORD_DEFAULT),
            'user_type' => '100',
            'nickname' => '创始人',
            'email' => 'admin@adminmine.com',
            'phone' => '16858888988',
            'signed' => '广阔天地，大有所为',
            'dashboard' => 'statistics',
            'created_by' => 0,
            'updated_by' => 0,
            'status' => 1,
            'created_at' => date('Y-m-d H:i:s'),
            'updated_at' => date('Y-m-d H:i:s'),
        ]);
        // 创建管理员角色
        $superRoleId = Db::table('system_role')->insertGetId([
            'name' => '超级管理员（创始人）',
            'code' => 'superAdmin',
            'data_scope' => 0,
            'sort' => 0,
            'created_by' => env('SUPER_ADMIN', 0),
            'updated_by' => 0,
            'status' => 1,
            'created_at' => date('Y-m-d H:i:s'),
            'updated_at' => date('Y-m-d H:i:s'),
            'remark' => '系统内置角色，不可删除',
        ]);
        Db::table('system_user_role')->insertGetId([
            'user_id' => $superAdminId,
            'role_id' => $superRoleId,
        ]);
        file_put_contents(BASE_PATH . '/.env', <<<ENV

# SUPER_ADMIN 超管id 必须要和数据表「system_user」保持一致 否则菜单会展示不全
SUPER_ADMIN = {$superAdminId}
# ADMIN_ROLE 超管角色id 必须要和数据表「system_role」保持一致 否则菜单会展示不全
ADMIN_ROLE = {$superRoleId}
ENV
            , FILE_APPEND);
    }

    protected function finish(): void
    {
        $i = 5;
        $this->output->write(PHP_EOL . $this->getGreenText('The installation is almost complete'), false);
        while ($i > 0) {
            $this->output->write($this->getGreenText('.'), false);
            --$i;
            sleep(1);
        }
        $this->line(PHP_EOL . sprintf('%s
MineAdmin Version: %s
default username: superAdmin
default password: admin123', $this->getInfo(), Mine::getVersion()), 'comment');
    }
}
