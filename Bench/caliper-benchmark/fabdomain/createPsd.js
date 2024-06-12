/*
* Licensed under the Apache License, Version 2.0 (the "License");
* you may not use this file except in compliance with the License.
* You may obtain a copy of the License at
*
* http://www.apache.org/licenses/LICENSE-2.0
*
* Unless required by applicable law or agreed to in writing, software
* distributed under the License is distributed on an "AS IS" BASIS,
* WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
* See the License for the specific language governing permissions and
* limitations under the License.
*/

'use strict';

const { WorkloadModuleBase } = require('@hyperledger/caliper-core');

// const colors = ['blue', 'red', 'green', 'yellow', 'black', 'purple', 'white', 'violet', 'indigo', 'brown'];
// const makes = ['Toyota', 'Ford', 'Hyundai', 'Volkswagen', 'Tesla', 'Peugeot', 'Chery', 'Fiat', 'Tata', 'Holden'];
// const models = ['Prius', 'Mustang', 'Tucson', 'Passat', 'S', '205', 'S22L', 'Punto', 'Nano', 'Barina'];
// const owners = ['Tomoko', 'Brad', 'Jin Soo', 'Max', 'Adrianna', 'Michel', 'Aarav', 'Pari', 'Valeria', 'Shotaro'];

const pkp = '02d37c24ddd0aaec7c3f1efe95d9cfbbec6b5c9c90291aac1ab1556a36d8cac97e'
const ct = '02d37c24ddd0aaec7c3f1efe95d9cfbbec6b5c9c90291aac1ab1556a36d8cac95a11111103069eccec991367ac7f65da48e8b72996fb85ee79f3c24a04c6c665a2b09370f6' 
const hpow = 'bcb607133f8f8b7622eeb54894fcef53c2ba51cf57138d3993b48dcaf39c1296'
const acc = 'e3dc3fc8df6b2af80428d69773bf37eec17b3b00017232df9f95578cb1c1c49d398755e87ff4a009db2c561d4bec55a1b31f21e88866832babff0b0c5b673d78e4713e613424624852cfd9558a3e2f4dbcb50f720a469e75a51e0d51d55e8cfe5a73e7b8fe5f82a1fb2b86df6317ad006166f21172b207b183fcbab6bad3e31a0b8f20e7aff6b3f7392856126d5e201c03'
const nonce = '1013edaec5b2b987d6a470b67a97e993e1fe67de18ac05e6acc7e6e321d9d1c9'

/**
 * Workload module for the benchmark round.
 */
class CreateCarWorkload extends WorkloadModuleBase {
    /**
     * Initializes the workload module instance.
     */
    constructor() {
        super();
        this.txIndex = 0;
    }

    /**
     * Assemble TXs for the round.
     * @return {Promise<TxStatus[]>}
     */
    async submitTransaction() {
        this.txIndex++;
        let pkp_psd = pkp + this.txIndex.toString();
        let ct_psd = ct + this.txIndex.toString();
        let hpow_psd = hpow + this.txIndex.toString();
        let acc_psd = acc + this.txIndex.toString();
        let nonce_psd = nonce + this.txIndex.toString();

        let args = {
            contractId: 'fabdomain',
            contractVersion: '1.0.1',
            contractFunction: 'createPsd',
            contractArguments: [pkp_psd, ct_psd, hpow_psd, acc_psd, nonce_psd],
            timeout: 30
        };

        await this.sutAdapter.sendRequests(args);
    }
}

/**
 * Create a new instance of the workload module.
 * @return {WorkloadModuleInterface}
 */
function createWorkloadModule() {
    return new CreateCarWorkload();
}

module.exports.createWorkloadModule = createWorkloadModule;
