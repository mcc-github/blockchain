/*
Copyright IBM Corp. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/
package chaincode.example;

import java.util.List;

import org.apache.commons.logging.Log;
import org.apache.commons.logging.LogFactory;
import org.mcc-github.blockchain.shim.ChaincodeBase;
import org.mcc-github.blockchain.shim.ChaincodeStub;

public class SimpleChaincode extends ChaincodeBase {

	private static Log _logger = LogFactory.getLog(SimpleChaincode.class);
	
	@Override
	public Response init(ChaincodeStub stub) {
		try {
			_logger.info("Init java simple chaincode");
			String func = stub.getFunction();
			if (!func.equals("init")) {
				return newErrorResponse("function other than init is not supported");
			}
			List<String> args = stub.getParameters();
			if (args.size() != 4) {
				newErrorResponse("Incorrect number of arguments. Expecting 4");
			}
			// Initialize the chaincode
			String account1Key = args.get(0);
			int account1Value = Integer.parseInt(args.get(1));
			String account2Key = args.get(2);
			int account2Value = Integer.parseInt(args.get(3));
			
			_logger.info(String.format("account %s, value = %s; account %s, value %s", account1Key, account1Value, account2Key, account2Value));
			stub.putState(account1Key, args.get(1));
			stub.putState(account2Key, args.get(3));
			
			return newSuccessResponse();
		} catch (Throwable e) {
			return newErrorResponse(e);
		}
	}

	@Override
	public Response invoke(ChaincodeStub stub) {
		try {
			_logger.info("Invoke java simple chaincode");
			String func = stub.getFunction();
			List<String> params = stub.getParameters();
			if (func.equals("invoke")) {
				return invoke(stub, params);
			}
			if (func.equals("delete")) {
				return delete(stub, params);
			}
			if (func.equals("query")) {
				return query(stub, params);
			}
			return newErrorResponse("Invalid invoke function name. Expecting one of: [\"invoke\", \"delete\", \"query\"]");
		} catch (Throwable e) {
			return newErrorResponse(e);
		}
	}

	private Response invoke(ChaincodeStub stub, List<String> args) {
		if (args.size() != 3) {
			return newErrorResponse("Incorrect number of arguments. Expecting 3");
		}
		String accountFromKey = args.get(0);
		String accountToKey = args.get(1);

		String accountFromValueStr = stub.getStateUTF8(accountFromKey);
		if (accountFromValueStr == null) {
			return newErrorResponse(String.format("Entity %s not found", accountFromKey));
		}
		int accountFromValue = Integer.parseInt(accountFromValueStr);

		String accountToValueStr = stub.getStateUTF8(accountToKey);
		if (accountToValueStr == null) {
			return newErrorResponse(String.format("Entity %s not found", accountToKey));
		}
		int accountToValue = Integer.parseInt(accountToValueStr);
		
		int amount = Integer.parseInt(args.get(2));

		if (amount > accountFromValue) {
			return newErrorResponse(String.format("not enough money in account %s", accountFromKey));
		}
		
		accountFromValue -= amount;
		accountToValue += amount;

		_logger.info(String.format("new value of A: %s", accountFromValue));
		_logger.info(String.format("new value of B: %s", accountToValue));
		
		stub.putState(accountFromKey, Integer.toString(accountFromValue));
		stub.putState(accountToKey, Integer.toString(accountToValue));
		
		_logger.info("Transfer complete");

		return newSuccessResponse("invoke finished successfully");
	}

	// Deletes an entity from state
	private Response delete(ChaincodeStub stub, List<String> args) {
		if (args.size() != 1) {
			return newErrorResponse("Incorrect number of arguments. Expecting 1");
		}
		String key = args.get(0);
		// Delete the key from the state in ledger
		stub.delState(key);
		return newSuccessResponse();
	}

	// query callback representing the query of a chaincode
	private Response query(ChaincodeStub stub, List<String> args) {
		if (args.size() != 1) {
			return newErrorResponse("Incorrect number of arguments. Expecting name of the person to query");
		}
		String key = args.get(0);
		//byte[] stateBytes
		String val	= stub.getStateUTF8(key);
		if (val == null) {
			return newErrorResponse(String.format("Error: state for %s is null", key));
		}
		_logger.info(String.format("Query Response:\nName: %s, Amount: %s\n", key, val));
		return newSuccessResponse(val);
	}
	
	public static void main(String[] args) {

		new SimpleChaincode().start(args);
	}
}
