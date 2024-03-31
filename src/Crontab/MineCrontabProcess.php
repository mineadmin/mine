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

namespace Mine\Crontab;

use Hyperf\Contract\StdoutLoggerInterface;
use Hyperf\Crontab\Event\CrontabDispatcherStarted;
use Hyperf\Crontab\Strategy\StrategyInterface;
use Hyperf\Di\Annotation\Inject;
use Hyperf\Process\AbstractProcess;
use Hyperf\Process\ProcessManager;
use Psr\Container\ContainerExceptionInterface;
use Psr\Container\ContainerInterface;
use Psr\Container\NotFoundExceptionInterface;
use Swoole\Coroutine;
use Swoole\Server;

class MineCrontabProcess extends AbstractProcess
{
    public string $name = 'MineAdmin Crontab';

    #[Inject]
    protected MineCrontabManage $mineCrontabManage;

    /**
     * @var Server
     */
    private $server;

    /**
     * @var MineCrontabScheduler
     */
    private $scheduler;

    /**
     * @var StrategyInterface
     */
    private $strategy;

    /**
     * @var StdoutLoggerInterface
     */
    private $logger;

    /**
     * @throws ContainerExceptionInterface
     * @throws NotFoundExceptionInterface
     */
    public function __construct(ContainerInterface $container)
    {
        parent::__construct($container);
        $this->scheduler = $container->get(MineCrontabScheduler::class);
        $this->strategy = $container->get(MineCrontabStrategy::class);
        $this->logger = $container->get(StdoutLoggerInterface::class);
    }

    public function bind($server): void
    {
        $this->server = $server;
        parent::bind($server);
    }

    /**
     * 是否自启进程.
     * @param Coroutine\Server|Server $server
     */
    public function isEnable($server): bool
    {
        return true;
    }

    /**
     * @throws ContainerExceptionInterface
     * @throws NotFoundExceptionInterface
     */
    public function handle(): void
    {
        $this->event->dispatch(new CrontabDispatcherStarted());
        while (ProcessManager::isRunning()) {
            $this->sleep();
            $crontabs = $this->scheduler->schedule();
            while (! $crontabs->isEmpty()) {
                /**
                 * @var MineCrontab $crontab
                 */
                $crontab = $crontabs->dequeue();
                $this->strategy->dispatch($crontab);
            }
        }
    }

    private function sleep()
    {
        $current = date('s', time());
        $sleep = 60 - $current;
        $this->logger->debug('MineAdmin Crontab dispatcher sleep ' . $sleep . 's.');
        $sleep > 0 && Coroutine::sleep($sleep);
    }
}
