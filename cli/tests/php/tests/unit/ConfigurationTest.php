<?php
declare(strict_types=1);
namespace Kapigen\Test\Unit;

use PHPUnit\Framework\TestCase;
use Kapigen\Configuration;

class ConfigurationTest extends TestCase
{
    public function testConfiguration(): void
    {
        $expectation = "test";
        $configuration = new Configuration($expectation);
        $this->assertEquals($expectation, $configuration->getTest());

        $expectation = "test2";
        $configuration->setTest($expectation);
        $this->assertEquals($expectation, $configuration->getTest());
    }
}
