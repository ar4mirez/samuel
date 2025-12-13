# Solidity Guide

> **Applies to**: Solidity 0.8+, Ethereum, EVM Chains, Smart Contracts, DeFi

---

## Core Principles

1. **Security First**: Smart contracts handle real value, bugs are expensive
2. **Gas Efficiency**: Every operation costs money
3. **Immutability**: Deployed code cannot be changed (plan for upgrades)
4. **Minimal Trust**: Assume all external calls are malicious
5. **Simplicity**: Less code = fewer bugs = lower attack surface

---

## Language-Specific Guardrails

### Solidity Version & Setup
- ✓ Use Solidity 0.8.x+ (built-in overflow checks)
- ✓ Lock pragma version: `pragma solidity 0.8.20;` (not `^0.8.20`)
- ✓ Use Hardhat or Foundry for development
- ✓ Use OpenZeppelin contracts for standard patterns
- ✓ Enable optimizer with appropriate runs (200 for normal, 1000+ for libraries)

### Code Style (Solidity Style Guide)
- ✓ Follow official Solidity Style Guide
- ✓ Use `camelCase` for functions and variables
- ✓ Use `PascalCase` for contracts, interfaces, structs, enums
- ✓ Use `SCREAMING_SNAKE_CASE` for constants and immutables
- ✓ Prefix internal/private functions with underscore: `_internalFunc`
- ✓ Prefix interfaces with `I`: `IERC20`
- ✓ Order: state variables, events, modifiers, constructor, functions

### Security (CRITICAL)
- ✓ Use Checks-Effects-Interactions pattern
- ✓ Use ReentrancyGuard for external calls
- ✓ Validate all inputs
- ✓ Use SafeERC20 for token transfers
- ✓ Avoid `tx.origin` for authentication
- ✓ Be careful with `delegatecall`
- ✓ Use access control (Ownable, AccessControl)
- ✓ Emit events for all state changes
- ✓ Get security audits before mainnet

### Gas Optimization
- ✓ Use `calldata` instead of `memory` for read-only function args
- ✓ Use `immutable` for constructor-set variables
- ✓ Use `constant` for compile-time constants
- ✓ Pack storage variables (smaller types together)
- ✓ Use `unchecked` blocks for safe arithmetic
- ✓ Cache storage reads in memory
- ✓ Use custom errors instead of require strings

### Visibility
- ✓ Use `external` for functions only called externally
- ✓ Use `public` for functions called internally and externally
- ✓ Use `internal` for functions used by derived contracts
- ✓ Use `private` for contract-specific functions
- ✓ Default to most restrictive visibility

---

## Project Structure

### Foundry Project
```
myproject/
├── foundry.toml
├── script/
│   └── Deploy.s.sol
├── src/
│   ├── MyContract.sol
│   ├── interfaces/
│   │   └── IMyContract.sol
│   └── libraries/
│       └── MyLibrary.sol
├── test/
│   ├── MyContract.t.sol
│   └── mocks/
│       └── MockERC20.sol
└── README.md
```

### Hardhat Project
```
myproject/
├── hardhat.config.ts
├── contracts/
│   ├── MyContract.sol
│   ├── interfaces/
│   └── libraries/
├── scripts/
│   └── deploy.ts
├── test/
│   └── MyContract.test.ts
├── package.json
└── README.md
```

### foundry.toml
```toml
[profile.default]
src = "src"
out = "out"
libs = ["lib"]
optimizer = true
optimizer_runs = 200
solc_version = "0.8.20"

[profile.default.fuzz]
runs = 256

[profile.ci.fuzz]
runs = 10000
```

---

## Security Patterns

### Checks-Effects-Interactions (CEI)
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

contract SecureWithdraw {
    mapping(address => uint256) public balances;

    // BAD: Vulnerable to reentrancy
    function withdrawBad() external {
        uint256 amount = balances[msg.sender];
        (bool success,) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
        balances[msg.sender] = 0;  // Effect after interaction!
    }

    // GOOD: CEI pattern
    function withdrawGood() external {
        // Checks
        uint256 amount = balances[msg.sender];
        require(amount > 0, "No balance");

        // Effects (before external call)
        balances[msg.sender] = 0;

        // Interactions (external call last)
        (bool success,) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }
}
```

### ReentrancyGuard
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract SecureVault is ReentrancyGuard {
    mapping(address => uint256) public balances;

    function withdraw(uint256 amount) external nonReentrant {
        require(balances[msg.sender] >= amount, "Insufficient balance");

        balances[msg.sender] -= amount;

        (bool success,) = msg.sender.call{value: amount}("");
        require(success, "Transfer failed");
    }
}
```

### Access Control
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

import "@openzeppelin/contracts/access/AccessControl.sol";

contract Treasury is AccessControl {
    bytes32 public constant ADMIN_ROLE = keccak256("ADMIN_ROLE");
    bytes32 public constant OPERATOR_ROLE = keccak256("OPERATOR_ROLE");

    constructor() {
        _grantRole(DEFAULT_ADMIN_ROLE, msg.sender);
        _grantRole(ADMIN_ROLE, msg.sender);
    }

    function emergencyWithdraw() external onlyRole(ADMIN_ROLE) {
        // Admin-only function
    }

    function processTransaction() external onlyRole(OPERATOR_ROLE) {
        // Operator function
    }
}
```

---

## Gas Optimization

### Storage Packing
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

// BAD: Uses 3 storage slots
contract BadPacking {
    uint256 a;     // Slot 0 (32 bytes)
    uint8 b;       // Slot 1 (1 byte, but takes full slot)
    uint256 c;     // Slot 2 (32 bytes)
    uint8 d;       // Slot 3 (1 byte, but takes full slot)
}

// GOOD: Uses 2 storage slots
contract GoodPacking {
    uint256 a;     // Slot 0 (32 bytes)
    uint256 c;     // Slot 1 (32 bytes)
    uint8 b;       // Slot 2 (1 byte)
    uint8 d;       // Slot 2 (1 byte, packed with b)
}
```

### Custom Errors
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

// BAD: Expensive string storage
contract BadErrors {
    function transfer(uint256 amount) external {
        require(amount > 0, "Amount must be greater than zero");
        require(balances[msg.sender] >= amount, "Insufficient balance");
    }
}

// GOOD: Custom errors (cheaper)
contract GoodErrors {
    error ZeroAmount();
    error InsufficientBalance(uint256 available, uint256 required);

    function transfer(uint256 amount) external {
        if (amount == 0) revert ZeroAmount();
        if (balances[msg.sender] < amount) {
            revert InsufficientBalance(balances[msg.sender], amount);
        }
    }
}
```

### Calldata vs Memory
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

contract DataLocation {
    // BAD: Copies array to memory
    function processBad(uint256[] memory data) external pure returns (uint256) {
        uint256 sum;
        for (uint256 i = 0; i < data.length; i++) {
            sum += data[i];
        }
        return sum;
    }

    // GOOD: Reads directly from calldata
    function processGood(uint256[] calldata data) external pure returns (uint256) {
        uint256 sum;
        for (uint256 i = 0; i < data.length; i++) {
            sum += data[i];
        }
        return sum;
    }
}
```

### Unchecked Blocks
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

contract UncheckedExample {
    // SAFE: i cannot overflow in practical loop bounds
    function sum(uint256[] calldata data) external pure returns (uint256 total) {
        uint256 length = data.length;
        for (uint256 i = 0; i < length;) {
            total += data[i];
            unchecked { ++i; }
        }
    }

    // SAFE: We already checked underflow condition
    function safeSub(uint256 a, uint256 b) external pure returns (uint256) {
        require(a >= b, "Underflow");
        unchecked {
            return a - b;
        }
    }
}
```

### Caching Storage Reads
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

contract StorageCache {
    uint256[] public items;

    // BAD: Reads storage in every iteration
    function sumBad() external view returns (uint256 total) {
        for (uint256 i = 0; i < items.length; i++) {
            total += items[i];
        }
    }

    // GOOD: Cache length in memory
    function sumGood() external view returns (uint256 total) {
        uint256 length = items.length;  // Cache
        for (uint256 i = 0; i < length; i++) {
            total += items[i];
        }
    }
}
```

---

## Common Patterns

### ERC20 Token
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

import "@openzeppelin/contracts/token/ERC20/ERC20.sol";
import "@openzeppelin/contracts/access/Ownable.sol";

contract MyToken is ERC20, Ownable {
    uint256 public constant MAX_SUPPLY = 1_000_000 * 10**18;

    constructor() ERC20("MyToken", "MTK") Ownable(msg.sender) {
        _mint(msg.sender, MAX_SUPPLY);
    }

    function burn(uint256 amount) external {
        _burn(msg.sender, amount);
    }
}
```

### Upgradeable Contract (UUPS)
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

import "@openzeppelin/contracts-upgradeable/proxy/utils/Initializable.sol";
import "@openzeppelin/contracts-upgradeable/proxy/utils/UUPSUpgradeable.sol";
import "@openzeppelin/contracts-upgradeable/access/OwnableUpgradeable.sol";

contract MyUpgradeable is Initializable, UUPSUpgradeable, OwnableUpgradeable {
    uint256 public value;

    /// @custom:oz-upgrades-unsafe-allow constructor
    constructor() {
        _disableInitializers();
    }

    function initialize(uint256 _value) external initializer {
        __Ownable_init(msg.sender);
        __UUPSUpgradeable_init();
        value = _value;
    }

    function setValue(uint256 _value) external onlyOwner {
        value = _value;
    }

    function _authorizeUpgrade(address newImplementation) internal override onlyOwner {}
}
```

### Staking Contract
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

import "@openzeppelin/contracts/token/ERC20/IERC20.sol";
import "@openzeppelin/contracts/token/ERC20/utils/SafeERC20.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol";

contract Staking is ReentrancyGuard {
    using SafeERC20 for IERC20;

    IERC20 public immutable stakingToken;
    IERC20 public immutable rewardToken;

    uint256 public rewardRate;
    uint256 public lastUpdateTime;
    uint256 public rewardPerTokenStored;

    mapping(address => uint256) public userRewardPerTokenPaid;
    mapping(address => uint256) public rewards;
    mapping(address => uint256) public balances;

    uint256 public totalSupply;

    error ZeroAmount();

    event Staked(address indexed user, uint256 amount);
    event Withdrawn(address indexed user, uint256 amount);
    event RewardPaid(address indexed user, uint256 reward);

    constructor(address _stakingToken, address _rewardToken, uint256 _rewardRate) {
        stakingToken = IERC20(_stakingToken);
        rewardToken = IERC20(_rewardToken);
        rewardRate = _rewardRate;
    }

    modifier updateReward(address account) {
        rewardPerTokenStored = rewardPerToken();
        lastUpdateTime = block.timestamp;
        if (account != address(0)) {
            rewards[account] = earned(account);
            userRewardPerTokenPaid[account] = rewardPerTokenStored;
        }
        _;
    }

    function rewardPerToken() public view returns (uint256) {
        if (totalSupply == 0) {
            return rewardPerTokenStored;
        }
        return rewardPerTokenStored +
            ((block.timestamp - lastUpdateTime) * rewardRate * 1e18) / totalSupply;
    }

    function earned(address account) public view returns (uint256) {
        return (balances[account] *
            (rewardPerToken() - userRewardPerTokenPaid[account])) / 1e18 +
            rewards[account];
    }

    function stake(uint256 amount) external nonReentrant updateReward(msg.sender) {
        if (amount == 0) revert ZeroAmount();

        totalSupply += amount;
        balances[msg.sender] += amount;

        stakingToken.safeTransferFrom(msg.sender, address(this), amount);

        emit Staked(msg.sender, amount);
    }

    function withdraw(uint256 amount) external nonReentrant updateReward(msg.sender) {
        if (amount == 0) revert ZeroAmount();

        totalSupply -= amount;
        balances[msg.sender] -= amount;

        stakingToken.safeTransfer(msg.sender, amount);

        emit Withdrawn(msg.sender, amount);
    }

    function getReward() external nonReentrant updateReward(msg.sender) {
        uint256 reward = rewards[msg.sender];
        if (reward > 0) {
            rewards[msg.sender] = 0;
            rewardToken.safeTransfer(msg.sender, reward);
            emit RewardPaid(msg.sender, reward);
        }
    }
}
```

---

## Testing

### Foundry Tests
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

import "forge-std/Test.sol";
import "../src/MyToken.sol";

contract MyTokenTest is Test {
    MyToken public token;
    address public owner = address(1);
    address public user = address(2);

    function setUp() public {
        vm.prank(owner);
        token = new MyToken();
    }

    function test_InitialSupply() public {
        assertEq(token.totalSupply(), 1_000_000 * 10**18);
        assertEq(token.balanceOf(owner), 1_000_000 * 10**18);
    }

    function test_Transfer() public {
        uint256 amount = 100 * 10**18;

        vm.prank(owner);
        token.transfer(user, amount);

        assertEq(token.balanceOf(user), amount);
        assertEq(token.balanceOf(owner), token.totalSupply() - amount);
    }

    function testFuzz_Transfer(uint256 amount) public {
        amount = bound(amount, 1, token.balanceOf(owner));

        vm.prank(owner);
        token.transfer(user, amount);

        assertEq(token.balanceOf(user), amount);
    }

    function test_RevertWhen_TransferExceedsBalance() public {
        uint256 amount = token.totalSupply() + 1;

        vm.prank(owner);
        vm.expectRevert();
        token.transfer(user, amount);
    }

    function test_Burn() public {
        uint256 burnAmount = 100 * 10**18;
        uint256 initialSupply = token.totalSupply();

        vm.prank(owner);
        token.burn(burnAmount);

        assertEq(token.totalSupply(), initialSupply - burnAmount);
    }
}
```

### Invariant Tests
```solidity
// SPDX-License-Identifier: MIT
pragma solidity 0.8.20;

import "forge-std/Test.sol";
import "../src/Vault.sol";

contract VaultInvariantTest is Test {
    Vault public vault;
    Handler public handler;

    function setUp() public {
        vault = new Vault();
        handler = new Handler(vault);

        targetContract(address(handler));
    }

    function invariant_TotalDepositsMustMatchBalance() public {
        assertEq(
            vault.totalDeposits(),
            address(vault).balance
        );
    }

    function invariant_UserBalancesSumToTotal() public {
        uint256 sum;
        address[] memory users = handler.getUsers();
        for (uint256 i = 0; i < users.length; i++) {
            sum += vault.balances(users[i]);
        }
        assertEq(sum, vault.totalDeposits());
    }
}

contract Handler is Test {
    Vault public vault;
    address[] public users;

    constructor(Vault _vault) {
        vault = _vault;
    }

    function deposit(uint256 amount) public {
        amount = bound(amount, 1, 10 ether);
        deal(msg.sender, amount);

        vm.prank(msg.sender);
        vault.deposit{value: amount}();

        _addUser(msg.sender);
    }

    function withdraw(uint256 amount) public {
        uint256 balance = vault.balances(msg.sender);
        amount = bound(amount, 0, balance);
        if (amount == 0) return;

        vm.prank(msg.sender);
        vault.withdraw(amount);
    }

    function getUsers() external view returns (address[] memory) {
        return users;
    }

    function _addUser(address user) internal {
        for (uint256 i = 0; i < users.length; i++) {
            if (users[i] == user) return;
        }
        users.push(user);
    }
}
```

---

## Tooling

### Foundry Commands
```bash
# Build
forge build

# Test
forge test
forge test -vvvv              # Verbose with traces
forge test --match-test testTransfer
forge test --gas-report

# Coverage
forge coverage
forge coverage --report lcov

# Deploy
forge script script/Deploy.s.sol --rpc-url $RPC_URL --broadcast

# Verify
forge verify-contract <address> MyContract --chain-id 1

# Format
forge fmt

# Gas snapshot
forge snapshot
forge snapshot --diff
```

### Hardhat Commands
```bash
# Compile
npx hardhat compile

# Test
npx hardhat test
npx hardhat test --grep "transfer"
npx hardhat coverage

# Deploy
npx hardhat run scripts/deploy.ts --network mainnet

# Verify
npx hardhat verify --network mainnet <address> <constructor-args>

# Console
npx hardhat console --network mainnet
```

### Security Tools
```bash
# Slither (static analysis)
slither .
slither . --print human-summary

# Mythril (symbolic execution)
myth analyze src/MyContract.sol

# Echidna (fuzzing)
echidna-test . --contract MyContract
```

---

## Common Pitfalls

### Don't Do This
```solidity
// Using tx.origin for auth
function withdraw() external {
    require(tx.origin == owner);  // Vulnerable to phishing!
}

// Unchecked external call success
function transfer(address to, uint256 amount) external {
    payable(to).transfer(amount);  // Can fail silently in some cases
    token.transfer(to, amount);     // Doesn't check return value!
}

// Floating pragma
pragma solidity ^0.8.0;  // Could compile with buggy version

// Timestamp dependence for randomness
function random() external view returns (uint256) {
    return uint256(keccak256(abi.encodePacked(block.timestamp)));
}

// Public by default
uint256 secretValue = 42;  // Actually public on blockchain!
```

### Do This Instead
```solidity
// Use msg.sender for auth
function withdraw() external {
    require(msg.sender == owner);
}

// Check call success
function transfer(address to, uint256 amount) external {
    (bool success,) = payable(to).call{value: amount}("");
    require(success, "Transfer failed");

    // Use SafeERC20
    IERC20(token).safeTransfer(to, amount);
}

// Lock pragma version
pragma solidity 0.8.20;

// Use Chainlink VRF for randomness
// Or commit-reveal scheme

// Everything on blockchain is public
// Use encryption off-chain, store hashes on-chain
```

---

## Security Checklist

Before deploying:

### Code Quality
- [ ] All functions have visibility specified
- [ ] Using SafeERC20 for token transfers
- [ ] Using ReentrancyGuard where needed
- [ ] Following CEI pattern
- [ ] Custom errors instead of require strings
- [ ] Events emitted for all state changes

### Access Control
- [ ] Functions have proper access modifiers
- [ ] Critical functions are protected
- [ ] Renouncing ownership is intentional

### Math & Logic
- [ ] Using Solidity 0.8+ (overflow protection)
- [ ] Division rounding handled correctly
- [ ] Edge cases tested (zero values, max values)

### External Interactions
- [ ] External call return values checked
- [ ] Callbacks are safe (reentrancy protected)
- [ ] Low-level calls have gas limits if needed

### Testing
- [ ] 100% test coverage on critical paths
- [ ] Fuzz tests for numeric inputs
- [ ] Invariant tests for key properties
- [ ] Fork tests against mainnet state

### Pre-Production
- [ ] Testnet deployment tested
- [ ] Gas optimized and estimated
- [ ] Security audit completed
- [ ] Bug bounty program ready

---

## References

- [Solidity Documentation](https://docs.soliditylang.org/)
- [OpenZeppelin Contracts](https://docs.openzeppelin.com/contracts/)
- [Foundry Book](https://book.getfoundry.sh/)
- [Solidity by Example](https://solidity-by-example.org/)
- [Smart Contract Security Best Practices](https://consensys.github.io/smart-contract-best-practices/)
- [SWC Registry](https://swcregistry.io/) (Smart Contract Weakness Classification)
- [Damn Vulnerable DeFi](https://www.damnvulnerabledefi.xyz/) (Security challenges)
- [Ethernaut](https://ethernaut.openzeppelin.com/) (Security wargame)
